package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steipete/foodoracli/internal/firebase"
)

func newSecretCmd(st *state) *cobra.Command {
	var printSecret bool
	var store bool
	var debug bool
	var country string

	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Fetch oauth client secret (from Firebase Remote Config)",
	}

	fetch := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch client secret for OAuth (client_id=android)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if country == "" {
				country = st.cfg.TargetCountryISO
			}
			if country == "" {
				country = "HU"
			}
			country = strings.ToUpper(country)

			rc := firebase.NewRemoteConfigClient(firebase.NetPincerHU)
			resp, err := rc.Fetch(cmd.Context())
			if err != nil {
				return err
			}

			raw, ok := resp.Entries["client_secrets"]
			if !ok || strings.TrimSpace(raw) == "" {
				if debug {
					fmt.Fprintf(cmd.ErrOrStderr(), "remote config keys: %s\n", strings.Join(mapKeys(resp.Entries), ", "))
				}
				return errors.New("remote config key client_secrets missing/empty")
			}

			var m map[string]string
			if err := json.Unmarshal([]byte(raw), &m); err != nil {
				return fmt.Errorf("client_secrets not JSON map: %w", err)
			}
			if debug {
				fmt.Fprintf(cmd.ErrOrStderr(), "client_secrets keys: %s\n", strings.Join(mapKeys(m), ", "))
			}

			rawByCountry := strings.TrimSpace(m[country])
			if rawByCountry == "" {
				return fmt.Errorf("client_secrets.%s missing/empty", country)
			}

			secret := ""
			// Newer configs: per-country JSON blob containing {android: "...", corp_android: "..."}.
			if strings.HasPrefix(rawByCountry, "{") {
				var per map[string]string
				if err := json.Unmarshal([]byte(rawByCountry), &per); err != nil {
					return fmt.Errorf("client_secrets.%s not JSON map: %w", country, err)
				}
				if debug {
					fmt.Fprintf(cmd.ErrOrStderr(), "client_secrets.%s keys: %s\n", country, strings.Join(mapKeys(per), ", "))
				}
				secret = strings.TrimSpace(per["android"])
			} else {
				// Older configs: per-country value is the secret itself.
				secret = rawByCountry
			}
			if secret == "" {
				return fmt.Errorf("client_secrets.%s.android missing/empty", country)
			}

			if store {
				st.cfg.ClientSecret = secret
				st.markDirty()
				fmt.Fprintln(cmd.OutOrStdout(), "stored")
				return nil
			}
			if printSecret {
				fmt.Fprintln(cmd.OutOrStdout(), secret)
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "ok (use --print to output, or --store to save into config)")
			return nil
		},
	}

	fetch.Flags().BoolVar(&printSecret, "print", false, "print secret to stdout")
	fetch.Flags().BoolVar(&store, "store", true, "store secret into config (default true)")
	fetch.Flags().BoolVar(&debug, "debug", false, "print debug info (keys only)")
	fetch.Flags().StringVar(&country, "country", "", "country code (defaults to config target ISO, else HU)")

	cmd.AddCommand(fetch)
	return cmd
}

func mapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
