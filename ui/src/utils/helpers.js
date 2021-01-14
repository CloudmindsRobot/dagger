import sshpk from 'sshpk'

export function copy(obj) {
  let newObj = {}
  for (let item in obj) {
    newObj[item] = obj[item]
  }
  return newObj
}

export function stringify(obj) {
  const items = []
  for (let item in obj) {
    items.push(item + '=' + obj[item])
  }
  return items.join('&')
}

export function clear(obj) {
  for (let item in obj) {
    obj[item] = null // 先让数据下沉到组件
    delete obj[item] // 再删属性
  }
}

export function merge(objA, objB) {
  let newObj = {}
  for (let item in objA) {
    newObj[item] = objA[item]
  }
  for (let item in objB) {
    newObj[item] = objB[item]
  }
  return newObj
}

export function sizeOfFmt(num) {
  let result = num
  const suffix = 'B'
  const units = ['', 'Ki', 'Mi', 'Gi', 'Ti', 'Pi', 'Ei', 'Zi']
  for (let index in units) {
    if (Math.abs(result) < 1024.0) {
      return `${result.toFixed(2)}${units[index]}${suffix}`
    }
    result /= 1024.0
  }
  return `${result.toFixed(2)}Yi${suffix}`
}

export function pythonDictToJson(pythonDict) {
  const items = []
  for (let pythonKey in pythonDict) {
    let newObj = {}
    newObj.name = pythonKey
    newObj.value = pythonDict[pythonKey]
    items.push(newObj)
  }
  return items
}

export function validateJWT(jwt) {
  if (!jwt) return false
  const user = JSON.parse(window.atob(jwt.split('.')[1]))
  const now = Date.parse(new Date()) / 1000
  return user.exp > now
}

export const rules = {
  requiredRules: [(v) => !!v || '该字段必填'],
  ipRules: [
    (v) =>
      !v ||
      /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-6])){3}$/.test(
        v,
      ) ||
      'IP地址格式不规范',
  ],
  ipListRules: [
    (v) => {
      if (!v) return true
      const ipList = v
        .replace(/^\s\s*/, '')
        .replace(/\s\s*$/, '')
        .split(/\s+/g) // 去掉首尾空格，并以任意空白作为分隔符进行分割
      for (let index in ipList) {
        if (!ipList[index]) return 'IP地址为空'
        if (
          !/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-6])){3}$/.test(
            ipList[index],
          )
        )
          return 'IP地址格式不规范'
      }
      return true
    },
  ],
  numberRules: [(v) => !v || /^[0-9]*$/.test(v) || '请输入整数'],
  emailRules: [
    (v) =>
      !v ||
      /^([A-Za-z0-9_\-.])+@([A-Za-z0-9_\-.])+\.([A-Za-z]{2,4})$/.test(v) ||
      '邮箱地址格式不规范',
  ],
  moboleLengthRules: [(v) => !v || v.length === 11 || '手机号码长度限制为11位'],
  passwordLengthRules: [(v) => !v || v.length >= 6 || '密码长度不能低于6'],
  noSubPathRules: [(v) => !v || !/\//.test(v) || '不能包含符号' / ''],
  sshPubKeyFingerprintRules: [
    (v) => {
      if (!v) return true
      try {
        sshpk
          .parseKey(v, 'ssh')
          .fingerprint('md5')
          .toString()
      } catch (err) {
        return 'ssh key 格式错误'
      }
      return true
    },
  ],
}

export function toLocaleString(val) {
  if (val) return val.toLocaleString()
  return val
}

export function getQueryString(name) {
  const reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i')
  const index = window.location.hash.indexOf('?')
  const r = window.location.hash.substr(index + 1).match(reg)
  if (r != null) {
    return unescape(r[2])
  }
  return null
}

export function formatDatetime(date, fmt) {
  const o = {
    'M+': date.getMonth() + 1,
    'd+': date.getDate(),
    'h+': date.getHours(),
    'm+': date.getMinutes(),
    's+': date.getSeconds(),
    'q+': Math.floor((date.getMonth() + 3) / 3),
    S: date.getMilliseconds(),
  }
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(
      RegExp.$1,
      (date.getFullYear() + '').substr(4 - RegExp.$1.length),
    )
  }
  for (var k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length == 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length),
      )
    }
  }
  return fmt
}

export function parserDatetime(date) {
  return new Date(Date.parse(date.replace(/-/g, '/')))
}
