gl2gh.exe
setlocal
set M="initial release"
call gi %M%
gh release create v0.0.1 --target main --notes %M%
