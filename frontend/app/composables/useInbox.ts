import { computed } from 'vue'
import { useApi, useMutation } from './useApi'
import type { Notification } from '~/types'

export const useInbox = () => {
  const {
    data: notifications,
    pending: isLoading,
    error,
    refresh: fetchInbox,
  } = useApi<Notification[]>('/api/v2/inbox', {
    key: 'inbox-notifications'
  })

  const markAsRead = async (id: string) => {
    try {
      await useMutation(`/api/v2/inbox/${id}/read`, {
        method: "PATCH",
      })
      if (notifications.value) {
        const index = notifications.value.findIndex(n => n.id === id)
        if (index !== -1) {
          notifications.value[index].is_read = true
        }
      }
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }

  const clearNotification = async (id: string) => {
    try {
      await useMutation(`/api/v2/inbox/${id}/clear`, {
        method: "PATCH",
      })
      if (notifications.value) {
        notifications.value = notifications.value.filter(n => n.id !== id)
      }
    } catch (error) {
      console.error('Failed to clear notification:', error)
    }
  }

  const clearAll = async () => {
    try {
      await useMutation(`/api/v2/inbox/clear-all`, {
        method: "PATCH",
      })
      if (notifications.value) {
        notifications.value = notifications.value.filter(n => !n.is_read)
      }
    } catch (error) {
      console.error('Failed to clear all read notifications:', error)
    }
  }

  const unreadCount = computed(() => {
    if (!notifications.value) return 0
    return notifications.value.filter(n => !n.is_read).length
  })

  return {
    notifications,
    isLoading,
    unreadCount,
    fetchInbox,
    markAsRead,
    clearNotification,
    clearAll
  }
}
