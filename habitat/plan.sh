pkg_name=bookend-scm-github
pkg_origin=stjohn
pkg_version="0.1.0"
pkg_source="https://github.com/stjohnjohnson/bookend-scm-github"
pkg_scaffolding=core/scaffolding-go
pkg_license=('BSD 3-clause')
pkg_deps=(core/git core/busybox-static)
pkg_bin_dirs=(bin)

scaffolding_go_build_deps=(github.com/fatih/color)

do_build() {
    export VERSION="${pkg_version}"
    scaffolding_go_build
}

do_install() {
    export VERSION="${pkg_version}"
    pushd "$scaffolding_go_src_path"
    make install
    popd
    cp -r "${scaffolding_go_gopath:?}/bin" "${pkg_prefix}/${bin}"

    wrap_bin "${pkg_prefix}/bin/bookend-scm-github"
}

wrap_bin() {
    local bin="$1"
    build_line "Adding wrapper $bin to ${bin}.real"
    mv -v "$bin" "${bin}.real"
    cat <<EOF > "${bin}"
#!$(pkg_path_for busybox-static)/bin/sh
set -e
if test -n "$DEBUG"; then set -x; fi
export GIT_PATH="$(pkg_path_for git)/bin/git"

exec ${bin}.real \$@
EOF
    chmod -v 755 "$bin"
}
