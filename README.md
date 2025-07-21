[![Go](https://img.shields.io/badge/code-Go-blue?logo=go)](README.md)
[![ES](https://img.shields.io/badge/lang-ES-red?logo=translate)](README.es.md)

# syspulse 📈

A terminal based system monitor written in Go using [Bubbletea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss) and [Gopsutil](https://github.com/shirou/gopsutil)

This is my capstone project for Boot.dev's backend bootcamp.
It shows live CPU and RAM gauges, along CPU, Memory, Running Processes and Disk stats with folder style tabs on a (Fallout New Vegas inspired) amber-themed TUI.

![syspulse demo](demo.gif)

## Why ?


I've been steadily learning programming and IT skills with the goal of transitioning my career into the fields of software development and IT operations.
Syspulse was born as a way to combine essential elements from both:

- 👨🏻‍💻From software development: I've applied Go's concurrency, data structures and modules along with a CLI/TUI interface implemented through external libraries such as [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

- 📈 From IT systems: I've integrated system resource monitoring through [Gopsutil](https://github.com/shirou/gopsutil?tab=readme-ov-file), learning how to track CPU, memory, running processes and disk data live.

## Tech Stack

- [Gopsutil](https://github.com/shirou/gopsutil?tab=readme-ov-file) - System metrics

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI Engine

- [Lipgloss](https://github.com/charmbracelet/lipgloss) - TUI Styling

- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI Components (tables, help messages)

## Description

This Go program, using bubbletea's ticker, periodically calls functions built using the gopsutil library to retrieve and update system metrics data.
This information is displayed on a terminal user interface constructed with charmbracelet _bubbletea_ library, styled with _lipgloss_ and using elements such as tables and togglable help messages from _bubbles_.

The content is displayed in different tabs corresponding to a category. The tabs and the information they display (in their respective order) are:

1. __CPU__
    - Total CPU percentual load gauge bar.
    - Table with time percentages the cpu spends on different operations.
2. __Memory__
    - Memory  percentual usage gauge bar.
    - Table with the amount of Total, Used, Free, Available, Buffer and Cached memory.
3. __Processes__
    - Table with the 7 top cpu demanding running processes
    - Including the process' ID, Name, Status, Runtime, Memory and CPU usage.
4. __Disks__
    - Table displaying the system's disk partitions.
    - Including the mountpoint, FsType, Total, Used and Free space

The metrics are gathered with functions from the __"systeminfo"__ module that uses _gopsutil_ library. They are periodically updated using a _bubbletea_ ticker with a modifiable time interval (currently set to 500 milliseconds).
If any error occurs during the data gathering process it is logged to __/logs/errors/systemstats.log__ using a logger function created with the _log/slog_ library.

_Note that there's a "config" folder and go file with no contents. It is intented to hold the module to be implemented to set system load thresholds to create performance logs and other configurations_

## Setup and Running Instructions (Linux/MacOS/WSL)

1. Prerequisites
    - Go (Golang) 1.19 or higher installed
    - Git (optional, for cloning the repository)

2. Clone or download the code
    ```bash
    git clone https://github.com/iegpeppino/syspulse.git
    cd syspulse
    ```
    Or [download the zip](https://github.com/iegpeppino/syspulse/refs/heads/main.zip) and extract it.

3. Initialize the Go module
    ```bash
    go mod init syspulse
    ```

4. Download dependencies
    ```bash
    go mod tidy
    ```

5. Run the app
    ```bash
    cd cmd
    go run .
    ```

6. Build an executable
    ```bash
    # While in cmd folder
    go build -o syspulse
    ./syspulse
    ```

7. Optional: Add to PATH
    ```bash
    sudo mv syspulse /usr/local/bin
    # Then, you can just use:
    syspulse
    ```
## Windows Instructions

1. Prerequisites
    - Go for Windows
    - Git for Windows (optional)
    - Windows terminal or PowerShell

2. Download or Clone
    ```cmd
    git clone https://github.com/iegpeppino/syspulse.git
    cd syspulse
    ```
     Or [download the zip](https://github.com/iegpeppino/syspulse/refs/heads/main.zip) and extract it.

3. Initialize and install dependencies
    ```cmd
    go mod init syspulse
    go mod tidy
    ```

4. Run program
    ```cmd
    cd cmd 
    go run .
    ```

5. Build (optional)
    - While on /cmd directory
    ```cmd
    go build -o syspulse.exe
    ./syspulse.exe

_Windows notes:_

_Its recommended to run the app from the terminal, and not double-clicking the executable._

_Ensure terminal supports ANSI escape sequences (use Windows Terminal or VSCode terminal)._

_Administrator permissions may be required to access some hardware metrics._

## Controls

- __(← / →) or (a / d)__ : Switch tabs

- __(q / ctrl + c / esc)__ : Quit

- __( h )__ : Show full help message


## 🤝 Contributing
### Submit a pull request
If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## Final notes

_Some other functions are to be implemented in the future_
    _such as testing network availability (pinging) and setting_
_system load thresholds to monitor and log reports when exceeded._

