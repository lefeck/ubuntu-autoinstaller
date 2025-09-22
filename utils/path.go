package utils

import (
	"fmt"
	"path/filepath"
)

type Path struct {
	RootDir string
}

// NewPath 构造一个 Path 实例
func NewPath(rootDir string) *Path {
	return &Path{RootDir: rootDir}
}

// 顶层目录
func (p *Path) BuildDir() string    { return filepath.Join(p.RootDir, "build") }
func (p *Path) DownloadDir() string { return filepath.Join(p.RootDir, "download") }
func (p *Path) ConfigDir() string   { return filepath.Join(p.RootDir, "config") }
func (p *Path) Boot() string        { return filepath.Join(p.RootDir, "BOOT") }

// build 下的细节目录
func (p *Path) Mount() string   { return filepath.Join(p.BuildDir(), "mnt") }
func (p *Path) BootISO() string { return filepath.Join(p.BuildDir(), "[BOOT]") }

func (p *Path) Packages() string { return filepath.Join(p.Mount(), "packages") }
func (p *Path) Scripts() string  { return filepath.Join(p.Mount(), "script") }

// 文件路径
func (p *Path) DownloadFile(fileName string) string {
	return filepath.Join(p.DownloadDir(), fileName)
}
func (p *Path) ScriptFile(fileName string) string {
	return filepath.Join(p.Scripts(), fileName)
}

// MD5SumFile
func (p *Path) MD5SumFile(md5sum string) string {
	return filepath.Join(p.BuildDir(), md5sum)
}

// grubConfigfile
func (p *Path) GrubConfigFile(grub string) string {
	return filepath.Join(p.BuildDir(), grub)
}

// loopBackConfigFile := gen.Path.LoopBackConfigFile(loopBack)
func (p *Path) LoopBackConfigFile(loopBack string) string {
	return filepath.Join(p.BuildDir(), loopBack)
}

// gen.Path.BuildPath(ISOhdpfxPath)
func (p *Path) ISOhdpfxFile(ISOhdpfx string) string {
	return filepath.Join(p.BuildDir(), ISOhdpfx)
}

// txtConfigPath
func (p *Path) TxtConfigFile(txt string) string {
	return filepath.Join(p.BuildDir(), txt)
}

func (p *Path) Sha256SumsFile(suffix string) string {
	return filepath.Join(p.DownloadDir(), fmt.Sprintf("SHA256SUMS-%s", suffix))
}

func (p *Path) Sha256SumsGPGFile(suffix string) string {
	return fmt.Sprintf("%s.gpg", p.Sha256SumsFile(suffix))
}

func (p *Path) KeyringFile(keyID string) string {
	return filepath.Join(p.DownloadDir(), fmt.Sprintf("%s.keyring", keyID))
}
func (p *Path) MetaDataFile(metaData string) string {
	return filepath.Join(p.BuildDir(), metaData)
}

func (p *Path) UserDataFile(userData string) string {
	return filepath.Join(p.BuildDir(), userData)
}
