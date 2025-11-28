module.exports = {
    apps: [
    {
        name: "MyServer",
        script: "./start_server.sh",
        instances: 4, // Number of instances to run
        autorestart: true,
        watch: false,
    },
    ],
};