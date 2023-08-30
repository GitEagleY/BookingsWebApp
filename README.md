<div align="center">
    <h1><strong>Bookings Web Application</strong></h1>
</div>

---

## 📋 Prerequisites

Before building and running the application, ensure you have the following in place:

- **Go**: The Go programming language is required for compiling and running the app.
- **PostgreSQL**: Make sure you have a PostgreSQL server up and running with a database.
- **Soda CLI**.
- **Mail server**: Recommended to use [mailhog](https://github.com/mailhog/MailHog).

---

#### [Go install instructions](https://go.dev/doc/install)

---

## 📦 Installing PostgreSQL

### For Ubuntu/Debian:

1.  Open a terminal.

2.  Update the package list and install PostgreSQL:

        sudo apt update
        sudo apt install postgresql postgresql-contrib

3.  Start and enable the PostgreSQL service:

        sudo systemctl start postgresql
        sudo systemctl enable postgresql

#### For non systemd systems:

         sudo service postgresql start
         sudo chkconfig postgresql on

### For macOS (using Homebrew):

1.  Open a terminal.

2.  Install PostgreSQL using Homebrew:

        brew install postgresql

3.  Start and enable the PostgreSQL service:

        brew services start postgresql

### For Windows:

1. Download the PostgreSQL installer from the official website: [PostgreSQL Downloads](https://www.postgresql.org/download/)

2. Run the installer and follow the on-screen instructions to complete the installation.

### For more detailed instructions see [PostgreSQL official website](https://www.postgresql.org/).

---

## 🚀 Build and Run Instructions

To successfully build and launch the application, please follow these steps:

1.  Begin by cloning or downloading the source code from the GitHub repository.

2.  Open your terminal or command prompt and navigate to the project directory.

3.  Execute the following command to conveniently download and install the necessary Go library dependencies:

        go mod download

4.  To manage database migrations, you'll need Soda CLI installed. You can install it using the following command:

        go install github.com/gobuffalo/pop/v5/soda@latest

5.  Perform database migrations using Soda CLI by running:

        soda reset

6.  You are now ready to build and run the application. Use the following command based on your operating system:

    For \*NIX:

        go build -o bookings cmd/web/*.go
        ./bookings

    For Windows:

        go build -o bookings.exe cmd/web/*.go
        ./bookings.exe

For mail sending - make sure you have a mail server running.

---

## 📜 License

This project is released under the BSD 3-Clause License.
