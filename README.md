
WELP: WeekEnd Lisp Project
==========================

Use [gimme][gimme_url] to fetch Go Modules-enabled Go version:

```
$ GIMME_GO_VERSION=1.12.4 gimme

unset GOOS;
unset GOARCH;
export GOROOT='/home/vytas/.gimme/versions/go1.12.4.linux.amd64';
export PATH="/home/vytas/.gimme/versions/go1.12.4.linux.amd64/bin:${PATH}";
go version >&2;

export GIMME_ENV="/home/vytas/.gimme/envs/go1.12.4.env"
```

[gimme_url]: https://github.com/travis-ci/gimme
