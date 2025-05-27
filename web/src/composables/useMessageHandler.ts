import { useMessage } from 'naive-ui'

/**
 * 消息处理结果类型
 */
export interface OperationResult {
  success: boolean
  message?: string
  error?: string
}

/**
 * 消息处理组合式函数
 */
export function useMessageHandler() {
  const message = useMessage()

  /**
   * 处理操作结果并显示相应消息
   */
  const handleResult = (result: OperationResult) => {
    if (result.success) {
      if (result.message) {
        message.success(result.message)
      }
    } else {
      if (result.error) {
        message.error(result.error)
      }
    }
  }

  /**
   * 处理异步操作
   */
  const handleAsyncOperation = async <T>(
    operation: () => Promise<T>,
    successMessage?: string,
    errorMessage?: string
  ): Promise<T | null> => {
    try {
      const result = await operation()
      if (successMessage) {
        message.success(successMessage)
      }
      return result
    } catch (error) {
      const errorMsg = errorMessage || (error instanceof Error ? error.message : '操作失败')
      message.error(errorMsg)
      return null
    }
  }

  return {
    message,
    handleResult,
    handleAsyncOperation
  }
}
