package util

import (
  "path"
  "os"
)

// ResolveImport takes a present working dir and a relative path and attemps to find the correct import
func ResolveImport(pwd, relativeImport string) string {
  // create fallback path in case TC_PACKAGES doesn't exist
  fallbackPath := path.Clean(path.Join(path.Clean(pwd), relativeImport))
  if fallbackStat, err := os.Stat(fallbackPath); err == nil && fallbackStat.IsDir() {
    fallbackPath = path.Join(fallbackPath, path.Base(fallbackPath) + ".tc")
  }

  // check the environment variable for a package import
  tcPackages := os.Getenv("TC_PACKAGES")
  if pkgStat, err := os.Stat(tcPackages); err == nil && pkgStat.IsDir() {
    absPath := path.Join(tcPackages, relativeImport)
    if absStat, err := os.Stat(absPath); err == nil && absStat.IsDir() {
      absPath = path.Join(absPath, path.Base(absPath) + ".tc")
    }
    if absFileStat, err := os.Stat(absPath); err == nil && !absFileStat.IsDir() {
      return absPath
    }
  }

  return fallbackPath
}
