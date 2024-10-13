The main purpose of this program is to implement local port forwarding through the SSH protocol, mapping a specified port (such as 8899) on the remote server to a local port, allowing you to access services on the remote server. The program continuously monitors the connection status and performs regular health checks (by running the `netstat` command on the remote server to check if port 8899 is in a listening state). If the connection is lost, it automatically retries after 5 seconds, ensuring stable operation of the forwarding service.

### Functionality:

1. **SSH Connection and Port Forwarding**:
   - The program first establishes an SSH connection with the remote server and maps a port on the remote server (e.g., 8899) to local port 22 through an SSH tunnel.
   
2. **Automatic Reconnection**:
   - If the connection is lost, the program automatically reconnects, avoiding manual intervention and ensuring service continuity.

3. **Health Check**:
   - Every 10 seconds, the program executes a remote `netstat` command via the SSH connection to check if port 8899 on the remote server is in a listening state, ensuring the forwarding is functioning normally.

4. **Bidirectional Data Forwarding**:
   - The program implements bidirectional data forwarding using `io.Copy`, ensuring smooth communication between remote and local ports.

### Using PM2 to Manage the Program

To facilitate long-term running of this program in a production environment, it is recommended to use PM2 to manage and monitor this Go program. PM2 is a powerful process management tool that can automatically start programs at system startup, monitor running status, and supports both Windows and macOS.

#### Installation and Usage Steps:

##### 1. Install PM2

PM2 is a Node.js-based tool, so you need to install Node.js first, then install PM2 via npm.

**Installation commands for Windows and macOS:**

```bash
# Install Node.js (if not installed)
# You can download it from the official Node.js website: https://nodejs.org/
# After installation, use the following command to install PM2

npm install pm2 -g
```

##### 2. Compile the Go Program

First, ensure that the Go environment is installed, then compile your Go program.

```bash
# Navigate to the program directory
cd ~/golang-unix-tools

# Compile the program
go build go_ssh_reverse_proxy.go
```

##### 3. Use PM2 to Manage the Compiled Go Program

Use the following commands to start and manage your Go program:

```bash
# Start the program and let it run in the background
pm2 start ./go_ssh_reverse_proxy --name "ssh-forward"

# View program status
pm2 status

# If you need to stop or restart
pm2 stop ssh-forward
pm2 restart ssh-forward
```

##### 4. Set PM2 to Start on Boot

You can set PM2 to automatically start your Go program when the system boots.

```bash
# Generate startup configuration
pm2 startup

# Save the current list of running processes
pm2 save
```

This way, whether on Windows or macOS, PM2 will automatically start and manage the Go program after system reboot.

### Supported Platforms:

- **Windows**: PM2 can be executed via command line or PowerShell, and the forwarding service can be automatically started at system boot.
- **macOS**: PM2 can be used through the terminal, and the service can be set to start automatically at boot. Download link: ![go_ssh_reverse_proxy_mac](./go_ssh_reverse_proxy_mac)

By following these steps, you can easily manage your SSH port forwarding program using PM2 on both Windows and macOS, ensuring stable access to remote server ports.
