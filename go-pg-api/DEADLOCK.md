## Working with deadlock in database
* Using [WaitGroup](https://pkg.go.dev/sync#example-WaitGroup)

Update in file `main.go`
```
// Endpoint for simulating the deadlock
e.POST("/transfer_deadlock", simulateDeadlockHandler)
```

Add function simulateDeadlockHandler
```

// simulateDeadlockHandler handles the API call to start the simulation.
func simulateDeadlockHandler(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Transaction 1: Lock account 1, then try to lock account 2
	go func() {
		defer wg.Done()
		err := runTransaction(1, 2, 50, "Transaction 1")
		if err != nil {
			log.Printf("Tx 1 FAILED: %v", err)
		}
	}()

	// Give Tx 1 a moment to acquire its first lock
	time.Sleep(100 * time.Millisecond)

	// Transaction 2: Lock account 2, then try to lock account 1
	go func() {
		defer wg.Done()
		err := runTransaction(2, 1, 50, "Transaction 2")
		if err != nil {
			log.Printf("Tx 2 FAILED: %v", err)
		}
	}()

	wg.Wait()
	return c.String(http.StatusOK, "Deadlock simulation started. Check logs for result.")
}

// runTransaction attempts to transfer amount from 'fromID' to 'toID'.
func runTransaction(fromID, toID, amount int, name string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to start transaction: %w", name, err)
	}
	defer tx.Rollback() // Rollback is safe to call even if Commit succeeds

	log.Printf("%s: Started. Acquiring lock on account %d.", name, fromID)

	// STEP 1: Lock 'from' account (Resource 1)
	// Use SELECT FOR UPDATE to acquire a row-level lock
	_, err = tx.Exec(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, fromID)
	if err != nil {
		return fmt.Errorf("%s: failed to lock account %d: %w", name, fromID, err)
	}

	// Introduce a delay here to ensure the other transaction can start and acquire its first lock,
	// creating the circular dependency (the heart of the deadlock).
	time.Sleep(500 * time.Millisecond)

	log.Printf("%s: Acquired lock on %d. Attempting to acquire lock on account %d.", name, fromID, toID)

	// STEP 2: Lock 'to' account (Resource 2)
	// This is where the deadlock occurs: Tx 1 waits for Tx 2's lock on account 2,
	// while Tx 2 (in its first step) is waiting for Tx 1's lock on account 1.
	_, err = tx.Exec(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, toID)
	if err != nil {
		// PostgreSQL detects the deadlock and one transaction will fail here with a '40P01' error.
		return fmt.Errorf("%s: **DEADLOCK POINT** failed to lock account %d: %w", name, toID, err)
	}

	// Perform the actual update/transfer (simplified)
	_, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		return fmt.Errorf("%s: failed to update balance for %d: %w", name, fromID, err)
	}

	_, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		return fmt.Errorf("%s: failed to update balance for %d: %w", name, toID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", name, err)
	}

	log.Printf("%s: SUCCESSFULLY committed transfer of %d from %d to %d.", name, amount, fromID, toID)
	return nil
}
```

List of URLs
* POST http://localhost:8080/transfer_deadlock

