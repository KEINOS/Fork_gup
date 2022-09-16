package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/gup/internal/print"
	"github.com/spf13/cobra"
)

func Test_list_not_found_go_command(t *testing.T) {
	t.Run("Not found go command", func(t *testing.T) {
		oldPATH := os.Getenv("PATH")
		if err := os.Setenv("PATH", ""); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Setenv("PATH", oldPATH); err != nil {
				t.Fatal(err)
			}
		}()

		orgStdout := print.Stdout
		orgStderr := print.Stderr
		pr, pw, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		print.Stdout = pw
		print.Stderr = pw

		if got := list(&cobra.Command{}, []string{}); got != 1 {
			t.Errorf("list() = %v, want %v", got, 1)
		}
		pw.Close()
		print.Stdout = orgStdout
		print.Stderr = orgStderr

		buf := bytes.Buffer{}
		_, err = io.Copy(&buf, pr)
		if err != nil {
			t.Error(err)
		}
		defer pr.Close()
		got := strings.Split(buf.String(), "\n")

		want := []string{
			`gup:ERROR: you didn't install golang: exec: "go": executable file not found in $PATH`,
			"",
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_list_gobin_is_empty(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name   string
		gobin  string
		args   args
		want   int
		stderr []string
	}{
		{
			name:  "gobin is empty",
			gobin: "./testdata/empty_dir",
			args:  args{},
			want:  1,
			stderr: []string{
				"gup:ERROR: unable to list up package: no package information",
				"",
			},
		},
		{
			name:  "$GOBIN is empty",
			gobin: "no_exist_dir",
			args:  args{},
			want:  1,
			stderr: []string{
				"gup:ERROR: can't get binary-paths installed by 'go install': open no_exist_dir: no such file or directory",
				"",
			},
		},
	}

	if err := os.Mkdir("./testdata/empty_dir", 0755); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldGoBin := os.Getenv("GOBIN")
			if err := os.Setenv("GOBIN", tt.gobin); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Setenv("GOBIN", oldGoBin); err != nil {
					t.Fatal(err)
				}
			}()

			orgStdout := print.Stdout
			orgStderr := print.Stderr
			pr, pw, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			print.Stdout = pw
			print.Stderr = pw

			if got := list(tt.args.cmd, tt.args.args); got != tt.want {
				t.Errorf("check() = %v, want %v", got, tt.want)
			}
			pw.Close()
			print.Stdout = orgStdout
			print.Stderr = orgStderr

			buf := bytes.Buffer{}
			_, err = io.Copy(&buf, pr)
			if err != nil {
				t.Error(err)
			}
			defer pr.Close()
			got := strings.Split(buf.String(), "\n")

			if diff := cmp.Diff(tt.stderr, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}

	err := os.Remove("./testdata/empty_dir")
	if err != nil {
		t.Fatal(err)
	}
}
