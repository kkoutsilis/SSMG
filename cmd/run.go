package cmd

import (
	"fmt"
	"io"

	"ssmg/internal"
)

func run(filePath, templatePath string, output io.Writer) error {
	config, err := internal.EmailConfigFromEnv()
	if err != nil {
		return err
	}

	participants, err := internal.LoadParticipants(filePath)
	if err != nil {
		return err
	}
	if len(participants) < 2 {
		return fmt.Errorf("at least two participants are required")
	}

	fmt.Fprintln(output, "Generating Secret Santa matches, HO HO HO...")
	matches := internal.GenerateMatches(participants)
	fmt.Fprintln(output, "Preparing and sending the emails, HO HO HO...")
	if err := internal.Send(matches, templatePath, config); err != nil {
		return err
	}
	fmt.Fprintln(output, "Happy gift hunting, and happy holidays!")
	return nil
}
