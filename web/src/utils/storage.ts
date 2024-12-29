const STORAGE_PREFIX = 'filenest_'

export const storage = {
  get(key: string) {
    const value = localStorage.getItem(STORAGE_PREFIX + key)
    return value ? JSON.parse(value) : null
  },

  set(key: string, value: any) {
    localStorage.setItem(STORAGE_PREFIX + key, JSON.stringify(value))
  },

  remove(key: string) {
    localStorage.removeItem(STORAGE_PREFIX + key)
  }
}
