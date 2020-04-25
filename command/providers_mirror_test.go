package command

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/mitchellh/cli"

	"github.com/hashicorp/terraform-svchost/disco"
)

func TestProvidersMirror(t *testing.T) {
	fixtureDir, err := filepath.Abs("testdata/providers-mirror")
	if err != nil {
		t.Fatal(err)
	}
	configDir := filepath.Join(fixtureDir, "test-config")
	defer testChdir(t, configDir)()

	if info, err := os.Stat("providers-mirror.tf"); err != nil || info.IsDir() {
		t.Fatalf("missing expected configuration fixture in %s", configDir)
	}

	outputDir := testTempDir(t)
	defer os.RemoveAll(outputDir)

	staticRegistryDir := filepath.Join(fixtureDir, "fake-registry")
	server := httptest.NewServer(http.FileServer(http.Dir(staticRegistryDir)))
	defer server.Close()
	registryBaseURL := server.URL + "/"
	t.Logf("Serving files from %s at %s", staticRegistryDir, registryBaseURL)

	services := disco.New()
	services.ForceHostServices(svchost.Hostname("example.com"), map[string]interface{}{
		"providers.v1": registryBaseURL,
	})

	ui := cli.NewMockUi()
	c := &ProvidersMirrorCommand{
		Meta: Meta{
			Ui:       ui,
			Services: services,
		},
	}

	status := c.Run([]string{"-platform=megadrive_m68k", "-platform=c64kernal_6510", outputDir})
	if want := 0; status != want {
		t.Fatalf("wrong response code %d; want %d\nstderr:\n%s", status, want, ui.ErrorWriter.String())
	}
}
