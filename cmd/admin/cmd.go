package admin

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "admin",
		Short: "AdminService commands",
		Long:  "Commands for backend.admin.v1.AdminService",
	}
}
