import React, { createContext, FC, useContext, useState } from 'react'

type MessageType = 'error' | 'warning' | 'info'

interface Message {
  text: string
  variant: MessageType
}

interface MessageContextType {
  message: Message
  showMessage: (text: string) => void
  showErrorMessage: (text: string) => void
}

const MessageContext = createContext<MessageContextType>({
  message: { text: '', variant: 'info' },
  showMessage: () => true,
  showErrorMessage: () => true,
})

const MessageProvider: FC = ({ children }) => {
  const [message, setMessage] = useState<Message>({ text: '', variant: 'info' })

  const showMessage = (text: string, variant: MessageType = 'info') => setMessage({ text, variant })
  const showErrorMessage = (text: string) => setMessage({ text, variant: 'error' })

  return (
    <MessageContext.Provider value={{ message, showMessage, showErrorMessage }}>{children}</MessageContext.Provider>
  )
}

export { MessageProvider }

export const useMessage = () => useContext(MessageContext)
