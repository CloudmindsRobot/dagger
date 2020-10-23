export function getCookie(name) {
  const reg = new RegExp('(^| )' + name + '=([^;]*)(;|$)')
  const arr = document.cookie.match(reg)
  if (arr) {
    return arr[2]
  }
  return null
}

// export function setCookie (cName, value, expiredays) {
//   const exdate = new Date()
//   exdate.setDate(exdate.getDate() + expiredays)
//   document.cookie = cName + '=' + escape(value) + ((expiredays === null) ? '' : ';expires=' + exdate.toGMTString())
// }

export function setCookie(cName, value, expireTimestamp) {
  let exdate
  if (expireTimestamp) {
    exdate = new Date(expireTimestamp)
  }
  document.cookie =
    cName +
    '=' +
    escape(value) +
    (!expireTimestamp ? '' : ';expires=' + exdate.toGMTString())
}

export function delCookie(name) {
  const exp = new Date()
  exp.setTime(exp.getTime() - 1)
  const cval = getCookie(name)
  if (cval !== null) {
    document.cookie = name + '=' + cval + ';expires=' + exp.toGMTString()
  }
}
