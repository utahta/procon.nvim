if exists('g:loaded_procon_nvim')
  finish
endif
let g:loaded_procon_nvim = 1

let s:plugin_name = 'procon'
let s:plugin_root = fnamemodify(resolve(expand('<sfile>:p')), ':h:h')
let s:plugin_cmd = [s:plugin_root . '/bin/' . s:plugin_name]

function! s:JobStart(host) abort
  return jobstart(s:plugin_cmd, { 'rpc': v:true })
endfunction

highlight link ProconAC String
highlight link ProconWA ErrorMsg

call remote#host#Register(s:plugin_name, '', function('s:JobStart'))
call remote#host#RegisterPlugin('procon', '0', [
\ {'type': 'command', 'name': 'ProconBuild', 'sync': 1, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'ProconCheck', 'sync': 1, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'ProconCurrentTask', 'sync': 1, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'ProconLogin', 'sync': 1, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'ProconReset', 'sync': 1, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'ProconSet', 'sync': 1, 'opts': {'nargs': '0'}},
\ {'type': 'command', 'name': 'ProconSubmit', 'sync': 1, 'opts': {'nargs': '0'}},
\ ])
