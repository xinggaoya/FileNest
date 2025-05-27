import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

// 配置dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

/**
 * 格式化日期时间
 * @param date 日期字符串或日期对象
 * @param format 格式化字符串，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的日期字符串
 */
export function formatDateTime(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(date).format(format)
}

/**
 * 格式化为相对时间
 * @param date 日期字符串或日期对象
 * @returns 相对时间字符串（如：2小时前、3天前）
 */
export function formatRelativeTime(date: string | Date): string {
  return dayjs(date).fromNow()
}

/**
 * 格式化为友好的时间显示
 * @param date 日期字符串或日期对象
 * @returns 友好的时间显示
 */
export function formatFriendlyTime(date: string | Date): string {
  const now = dayjs()
  const target = dayjs(date)
  const diffInDays = now.diff(target, 'day')

  if (diffInDays === 0) {
    // 今天
    return target.format('HH:mm')
  } else if (diffInDays === 1) {
    // 昨天
    return '昨天 ' + target.format('HH:mm')
  } else if (diffInDays < 7) {
    // 一周内
    return target.format('dddd HH:mm')
  } else if (target.year() === now.year()) {
    // 今年
    return target.format('MM-DD HH:mm')
  } else {
    // 其他年份
    return target.format('YYYY-MM-DD')
  }
}

/**
 * 检查日期是否为今天
 * @param date 日期字符串或日期对象
 * @returns 是否为今天
 */
export function isToday(date: string | Date): boolean {
  return dayjs(date).isSame(dayjs(), 'day')
}

/**
 * 检查日期是否为昨天
 * @param date 日期字符串或日期对象
 * @returns 是否为昨天
 */
export function isYesterday(date: string | Date): boolean {
  return dayjs(date).isSame(dayjs().subtract(1, 'day'), 'day')
}

/**
 * 检查日期是否为本周
 * @param date 日期字符串或日期对象
 * @returns 是否为本周
 */
export function isThisWeek(date: string | Date): boolean {
  return dayjs(date).isSame(dayjs(), 'week')
}
