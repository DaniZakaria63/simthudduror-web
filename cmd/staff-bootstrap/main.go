// Command staff-bootstrap creates the first staff account through an
// operator-only local action. There is deliberately no HTTP bootstrap route.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"maulid.id/internal/account"
	"maulid.id/internal/platform/config"
	"maulid.id/internal/platform/database"
)

func main() {
	email := flag.String("email", "", "staff email address")
	name := flag.String("name", "", "staff display name")
	passwordStdin := flag.Bool("password-stdin", false, "read the password from the first stdin line (for operators)")
	flag.Parse()

	if strings.TrimSpace(*email) == "" || strings.TrimSpace(*name) == "" {
		fmt.Fprintln(os.Stderr, "usage: staff-bootstrap --email <email> --name <name> [--password-stdin]")
		os.Exit(2)
	}
	password, err := readPassword(*passwordStdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "password:", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration:", err)
		os.Exit(1)
	}
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "database:", err)
		os.Exit(1)
	}
	defer db.Close()

	service := account.NewService(account.NewRepository(db), 24*time.Hour)
	created, err := service.BootstrapStaff(context.Background(), *email, *name, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Akun staf dibuat: %s (%s)\n", created.Email, created.ID)
}

func readPassword(fromStdin bool) (string, error) {
	if fromStdin {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			return "", fmt.Errorf("tidak ada kata sandi pada stdin")
		}
		password := strings.TrimRight(line, "\r\n")
		if password == "" {
			return "", fmt.Errorf("kata sandi tidak boleh kosong")
		}
		return password, nil
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", fmt.Errorf("stdin bukan terminal; gunakan --password-stdin")
	}
	fmt.Fprint(os.Stderr, "Kata sandi staf: ")
	first, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	fmt.Fprint(os.Stderr, "Ulangi kata sandi: ")
	second, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	if string(first) != string(second) {
		return "", fmt.Errorf("kata sandi tidak sama")
	}
	return string(first), nil
}
