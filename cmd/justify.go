/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/rightfoot-consulting/txttool/textp"
	"github.com/spf13/cobra"
)

// justifyCmd represents the justify command
var justifyCmd = &cobra.Command{
	Use:   "justify",
	Short: "Justify the lines in a text file to adhere to specific column width",
	Long: `This command will read in a text file and perform the following operations:
	   1. Remove all duplicate spaces in each line as well as any white space from empty lines.
	   2. Remove any extraneous line breaks, defined as any line break not followed by an empty line.
	   3. Tokenize each word or specified phrase and then justify the text by breaking up the lines
	   to fit a specific line width.
	   4. Write the file to an output file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("justify called")
		input, err := cmd.Flags().GetString("in")
		if err != nil {
			log.Fatalf("error getting --in flag: %v", err)
		}
		output, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatalf("error getting --out flag: %v", err)
		}
		len, err := cmd.Flags().GetInt("len")
		if err != nil {
			log.Fatalf("error getting --len flag: %v", err)
		}
		err = textp.JustifyText(input, output, len)
		if err != nil {
			log.Fatalf("error in processing: %v", err)
		}
		log.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(justifyCmd)
	justifyCmd.Flags().StringP("in", "i", "", "Input file to process")
	justifyCmd.Flags().StringP("out", "o", "", "File to write the result to, if unspecified, the input filename plust a '.txt' extension will be used.")
	justifyCmd.Flags().IntP("len", "l", 120, "Target Line length. Defaults to 120")
}
