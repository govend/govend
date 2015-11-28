package commands

var (
	importsPriority string
	importsDisplay  bool
	importsErrors   bool
	importsList     bool
	importsWrite    bool
)

/*

func init() {
	ImportsCMD.Flags().StringVarP(&importsPriority, "priority", "p", "", "Comma seperated import path prefixes that if matched, take priority")
	ImportsCMD.Flags().BoolVarP(&importsDisplay, "display", "d", false, "Display diffs instead of rewriting files")
	ImportsCMD.Flags().BoolVarP(&importsErrors, "errors", "e", false, "Report all errors (not just the first 10 on different lines)")
	ImportsCMD.Flags().BoolVarP(&importsList, "list", "l", false, "List files whose formatting differs from goimport's")
	ImportsCMD.Flags().BoolVarP(&importsWrite, "write", "w", false, "Write result to (source) file instead of stdout")
}

// ImportsCMD describes the imports command.
var ImportsCMD = &cobra.Command{
	Use:   "imports",
	Short: "Rewrite import paths, prioritizing the project vendor directory.",
	Long:  "Use this command to for goimports functionality that prioritizes the project vendor directory imports.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 1 {
			if err := packages.ImportsCMD(importsDisplay, importsErrors, importsList, importsWrite, importsPriority, vendorDir, args); err != nil {
				log.Fatal(err)
			}
		}

		if err := packages.ImportsCMD(importsDisplay, importsErrors, importsList, importsWrite, importsPriority, vendorDir, []string{}); err != nil {
			log.Fatal(err)
		}
	},
}
*/
