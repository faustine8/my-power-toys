package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/faustine8/my-power-toys/internal/config"
	"github.com/faustine8/my-power-toys/internal/opener"
	"github.com/faustine8/my-power-toys/internal/project"
	"github.com/spf13/cobra"
)

// Version is set at build time or defaults to "dev".
var Version = "dev"

type rootOptions struct {
	StorePath   string
	Getwd       func() (string, error)
	RunOpenCode func(string) error
}

func newRootCommand(options rootOptions) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "po",
		Short:         "po - A short-lived project opener for developers.",
		Long:          "po is a local development productivity tool.\nIt manages project directories and opens them quickly from the shell.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(newVersionCommand())
	rootCmd.AddCommand(newAddCommand(options))
	rootCmd.AddCommand(newListCommand(options))
	rootCmd.AddCommand(newRemoveCommand(options))
	rootCmd.AddCommand(newPickCommand(options))
	rootCmd.AddCommand(newOpenCodeCommand(options))

	return rootCmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of po",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "po "+Version)
		},
	}
}

func newAddCommand(options rootOptions) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Register the current directory as a project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := getwd(options)
			if err != nil {
				return err
			}

			store, err := storeFor(options)
			if err != nil {
				return err
			}
			file, err := store.Load()
			if err != nil {
				return err
			}

			updated, added, err := project.Add(file, cwd, name)
			if err != nil {
				return err
			}
			if err := store.Save(updated); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "added\t%s\t%s\n", added.Name, added.Path)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "project alias")
	return cmd
}

func newListCommand(options rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List registered projects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storeFor(options)
			if err != nil {
				return err
			}
			file, err := store.Load()
			if err != nil {
				return err
			}

			for _, item := range project.List(file) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", item.Name, item.Path)
			}
			return nil
		},
	}
}

func newRemoveCommand(options rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a registered project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storeFor(options)
			if err != nil {
				return err
			}
			file, err := store.Load()
			if err != nil {
				return err
			}

			updated, removed, err := project.Remove(file, args[0])
			if err != nil {
				return err
			}
			if err := store.Save(updated); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "removed\t%s\n", removed.Name)
			return nil
		},
	}
}

func newPickCommand(options rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "pick [query]",
		Short: "Select or search a project and print its path",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) == 1 {
				query = args[0]
			}

			selected, ok, err := pickProject(cmd, options, query)
			if err != nil {
				return err
			}
			if ok {
				fmt.Fprintln(cmd.OutOrStdout(), selected.Path)
			}
			return nil
		},
	}
}

func pickProject(cmd *cobra.Command, options rootOptions, query string) (config.Project, bool, error) {
	projects, err := loadProjects(options)
	if err != nil {
		return config.Project{}, false, err
	}
	if strings.TrimSpace(query) != "" {
		projects = project.Search(projects, query)
		if len(projects) == 0 {
			return config.Project{}, false, fmt.Errorf("no project matches query: %s", query)
		}
		if len(projects) == 1 {
			return projects[0], true, nil
		}
	}
	return project.Select(projects, cmd.InOrStdin(), cmd.ErrOrStderr())
}

func newOpenCodeCommand(options rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "oc",
		Short: "Select a project and run opencode in that directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			selected, ok, err := selectProject(cmd, options)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			return runOpenCode(options, selected.Path)
		},
	}
}

func selectProject(cmd *cobra.Command, options rootOptions) (config.Project, bool, error) {
	projects, err := loadProjects(options)
	if err != nil {
		return config.Project{}, false, err
	}
	return project.Select(projects, cmd.InOrStdin(), cmd.ErrOrStderr())
}

func loadProjects(options rootOptions) ([]config.Project, error) {
	store, err := storeFor(options)
	if err != nil {
		return nil, err
	}
	file, err := store.Load()
	if err != nil {
		return nil, err
	}
	return project.List(file), nil
}

func storeFor(options rootOptions) (config.Store, error) {
	path := options.StorePath
	if path == "" {
		var err error
		path, err = config.DefaultPath()
		if err != nil {
			return config.Store{}, err
		}
	}
	return config.Store{Path: path}, nil
}

func getwd(options rootOptions) (string, error) {
	if options.Getwd != nil {
		return options.Getwd()
	}
	return os.Getwd()
}

func runOpenCode(options rootOptions, dir string) error {
	if options.RunOpenCode != nil {
		return options.RunOpenCode(dir)
	}
	return opener.DefaultRunner().RunOpenCode(dir)
}

// Execute runs the root command.
func Execute() {
	rootCmd := newRootCommand(rootOptions{})
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
