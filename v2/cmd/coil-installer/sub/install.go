package sub

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func installCniConf(cniConfName, cniEtcDir, cniNetConf, cniNetConfFile string) error {
	data := []byte(cniNetConf)
	if cniNetConf == "" {
		bData, err := ioutil.ReadFile(cniNetConfFile)
		if err != nil {
			return err
		}
		data = bData
	}

	err := os.MkdirAll(cniEtcDir, 0755)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(cniEtcDir)
	if err != nil {
		return err
	}
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		err := os.Remove(filepath.Join(cniEtcDir, fi.Name()))
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filepath.Join(cniEtcDir, cniConfName))
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Chmod(0644)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}

func installCoil(coilPath, cniBinDir string) error {
	f, err := os.Open(coilPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = os.MkdirAll(cniBinDir, 0755)
	if err != nil {
		return err
	}

	g, err := ioutil.TempFile(cniBinDir, ".tmp")
	if err != nil {
		return err
	}
	defer func() {
		g.Close()
		os.Remove(g.Name())
	}()

	_, err = io.Copy(g, f)
	if err != nil {
		return err
	}

	err = g.Chmod(0755)
	if err != nil {
		return err
	}

	err = g.Sync()
	if err != nil {
		return err
	}

	return os.Rename(g.Name(), filepath.Join(cniBinDir, "coil"))
}
