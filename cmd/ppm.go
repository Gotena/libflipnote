package cmd

import (
	"fmt"
	"image"
	"image/gif"
	"os"
	"strings"

	"github.com/Clinet/ffgoconv"
	"github.com/Gotena/libflipnote/ppm"
	"github.com/Gotena/libflipnote/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var decode = &cobra.Command{
	Use:   "ppmdecode",
	Short: "decode a .ppm file",
	Run:   decodeFlipnote,
}

var decodeMeta = &cobra.Command{
	Use:   "ppmmeta",
	Short: "print metadata about a .ppm file",
}

var decodeSound = &cobra.Command{
	Use:   "ppmsound",
	Short: "decode a .ppm file's audio data",
}

func initPpmCmds() {
	decode.Flags().StringP("input", "i", "input.ppm", "input file.ppm")
	decode.Flags().StringP("output", "o", "output.mp4", "output file.mp4")

	decode.MarkFlagRequired("input")
	decode.MarkFlagRequired("output")

	rootCmd.AddCommand(decode)
	rootCmd.AddCommand(decodeMeta)
	rootCmd.AddCommand(decodeSound)
}

func decodeFlipnote(cmd *cobra.Command, args []string) {
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		cmd.Usage()
		return
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		cmd.Usage()
		return
	}

	if !strings.HasSuffix(input, ".ppm") {
		cmd.Usage()
		return
	}

	name := strings.TrimSuffix(input, ".ppm")
	nameGIF := fmt.Sprintf("%s.gif", name)
	nameWAV := fmt.Sprintf("%s.wav", name)

	os.Remove(nameGIF)
	os.Remove(nameWAV)
	os.Remove(output)

	p := progressbar.New(100)

	flip, err := ppm.ReadFile(input)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		os.Exit(1)
	}
	images := make([]*image.Paletted, flip.FrameCount)
	for i := uint16(0); i < flip.FrameCount; i++ {
		images[i] = flip.Frames[i].GetImage()
	}

	gifFile, err := os.Create(nameGIF)
	if err != nil {
		fmt.Printf("error creating file: %v\n", err)
		os.Exit(1)
	}

	timings := make([]int, flip.FrameCount)

	gif.EncodeAll(gifFile, &gif.GIF{
		Image: images,
		Delay: timings,
	})

	utils.IncrementBar(p, 25)

	err = gifFile.Close()
	if err != nil {
		fmt.Printf("error closing file: %v\n", err)
		os.Exit(1)
	}

	wavFile, err := os.Create(nameWAV)
	if err != nil {
		fmt.Printf("error creating file: %v\n", err)
		os.Exit(1)
	}

	err = flip.Audio.Export(wavFile, flip, 32768)
	if err != nil {
		fmt.Printf("error exporting audio: %v\n", err)
		os.Exit(1)
	}

	err = wavFile.Close()
	if err != nil {
		fmt.Printf("error closing file: %v\n", err)
		os.Exit(1)
	}

	utils.IncrementBar(p, 25)

	ffmpeg, err := ffgoconv.NewFFmpeg(name, []string{"-hide_banner", "-stats",
		"-r", fmt.Sprintf("%.1f", flip.Framerate),
		"-hwaccel", "auto",
		"-i", nameGIF,
		"-i", nameWAV,
		output,
		"-pix_fmt", "yuv420p",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-threads", "24",
	})
	if err != nil {
		fmt.Printf("error creating ffmpeg: %v\n", err)
		os.Exit(1)
	}

	if err := ffmpeg.Run(); err != nil {
		fmt.Printf("error running ffmpeg: %v\n", err)
		os.Exit(1)
	}

	utils.IncrementBarSync(p, 50)

	os.Remove(nameGIF)
	os.Remove(nameWAV)
}
