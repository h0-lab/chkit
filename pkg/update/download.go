package update

import (
	"io"

	"fmt"

	"runtime"

	"net/http"

	"io/ioutil"
	"os"
	"path"

	"github.com/blang/semver"
	"github.com/cheggaaa/pb"
	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"gopkg.in/urfave/cli.v2"
)

type LatestChecker interface {
	LatestVersion() (semver.Version, error)
}

type LatestDownloader interface {
	LatestDownload() (io.ReadCloser, error)
}

type LatestCheckerDownloader interface {
	LatestChecker
	LatestDownloader
}

const (
	ErrUpdateCheck    = chkitErrors.Err("unable to check latest version")
	ErrUpdateDownload = chkitErrors.Err("unable to download latest version")
)

const (
	MaxFileSize = 50 * (1 << 20) // 50 megabytes
)

func DownloadFileName(version semver.Version) string {
	extension := "tar.gz"
	if runtime.GOOS == "windows" {
		extension = "zip"
	}
	return fmt.Sprintf("chkit_%s_%s_v%s.%s", runtime.GOOS, runtime.GOARCH, version, extension)
}

type GithubLatestCheckerDownloader struct {
	client      *resty.Client
	log         *logrus.Logger
	ctx         *cli.Context
	downloadUrl string
}

func NewGithubLatestCheckerDownloader(ctx *cli.Context, owner, repo string) *GithubLatestCheckerDownloader {
	return &GithubLatestCheckerDownloader{
		ctx: ctx,
		log: util.GetLog(ctx),
		client: resty.New().
			SetHostURL(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)).
			SetDebug(true).
			SetLogger(util.GetLog(ctx).WriterLevel(logrus.DebugLevel)),
		downloadUrl: fmt.Sprintf("https://github.com/%s/%s/releases/download", owner, repo),
	}
}

func (gh *GithubLatestCheckerDownloader) LatestVersion() (semver.Version, error) {
	gh.log.Debug("get latest version from github")

	var latestVersionResp struct {
		LatestVersion semver.Version `json:"tag_name"`
	}

	_, err := gh.client.R().SetResult(&latestVersionResp).Get("/latest")
	if err != nil {
		return semver.MustParse("0.0.1-alpha"), chkitErrors.Wrap(ErrUpdateCheck, err)
	}

	return latestVersionResp.LatestVersion, nil
}

func (gh *GithubLatestCheckerDownloader) LatestDownload() (io.ReadCloser, error) {
	gh.log.Debug("download update")

	latestVersion, err := gh.LatestVersion()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", gh.downloadUrl, latestVersion, DownloadFileName(latestVersion)))
	if err != nil {
		return nil, chkitErrors.Wrap(ErrUpdateDownload, err)
	}

	bar := pb.New64(resp.ContentLength)

	return bar.NewProxyReader(resp.Body), nil
}

type FileSystemLatestCheckerDownloader struct {
	baseDir string
	log     *logrus.Logger
	ctx     *cli.Context
}

func NewFileSystemLatestCheckerDownloader(ctx *cli.Context, baseDir string) *FileSystemLatestCheckerDownloader {
	return &FileSystemLatestCheckerDownloader{
		baseDir: baseDir,
		ctx:     ctx,
		log:     util.GetLog(ctx),
	}
}

func (fs *FileSystemLatestCheckerDownloader) LatestVersion() (semver.Version, error) {
	fs.log.Debug("get latest version from filesystem")

	verFile, err := os.Open(path.Join(fs.baseDir, "version"))
	if err != nil {
		return semver.MustParse("0.0.1-alpha"), chkitErrors.Wrap(ErrUpdateCheck, err)
	}
	defer verFile.Close()

	ver, err := ioutil.ReadAll(verFile)
	if err != nil {
		return semver.MustParse("0.0.1-alpha"), chkitErrors.Wrap(ErrUpdateCheck, err)
	}

	sv, err := semver.ParseTolerant(string(ver))
	if err != nil {
		return semver.MustParse("0.0.1-alpha"), chkitErrors.Wrap(ErrUpdateCheck, err)
	}

	return sv, nil
}

func (fs *FileSystemLatestCheckerDownloader) LatestDownload() (io.ReadCloser, error) {
	fs.log.Debug("get latest version package from filesystem")

	latestVersion, err := fs.LatestVersion()
	if err != nil {
		return nil, err
	}

	pkg, err := os.Open(path.Join(fs.baseDir, DownloadFileName(latestVersion)))
	if err != nil {
		return nil, chkitErrors.Wrap(ErrUpdateDownload, err)
	}

	return pkg, nil
}