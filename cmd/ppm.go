package cmd

import "github.com/spf13/cobra"

var decode = &cobra.Command{
	Use:   "ppmdecode",
	Short: "decode a .ppm file",
}

var decodeMeta = &cobra.Command{
	Use:   "ppmmeta",
	Short: "print metadata about a .ppm file",
}

var decodeSound = &cobra.Command{
	Use:   "ppmsound",
	Short: "decode a .ppm file's audio data",
}
