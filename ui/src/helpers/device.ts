const isMobile = {
  Android() {
    return navigator.userAgent.match(/Android/i)
  },

  BlackBerry() {
    return navigator.userAgent.match(/BlackBerry/i)
  },

  iOS() {
    return navigator.userAgent.match(/iPhone|iPad|iPod/i)
  },

  Opera() {
    return navigator.userAgent.match(/Opera Mini/i)
  },

  Windows() {
    return navigator.userAgent.match(/IEMobile/i)
  },

  any() {
    return this.Android() || this.BlackBerry() || this.iOS() || this.Opera() || this.Windows()
  },
}

export const isMobileDevice = () => isMobile.any()

export const isInstalled = () => window.matchMedia('(display-mode: standalone)').matches
