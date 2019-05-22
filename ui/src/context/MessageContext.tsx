import React, { createContext, ReactNode, useState } from 'react'

interface Message {
  text: string
  isError: boolean
}

interface MessageContextType {
  message: Message
  showMessage: (text: string) => void
  showErrorMessage: (text: string) => void
}

const MessageContext = createContext<MessageContextType>({
  message: { text: '', isError: false },
  showMessage: () => true,
  showErrorMessage: () => true
})

interface Props {
  children: ReactNode
}

const MessageProvider = ({ children }: Props) => {
  const [message, setMessage] = useState<Message>({ text: '', isError: false })

  const showMessage = (text: string) => setMessage({ text, isError: false })
  const showErrorMessage = (text: string) => setMessage({ text, isError: true })

  return (
    <MessageContext.Provider value={{ message, showMessage, showErrorMessage }}>{children}</MessageContext.Provider>
  )
}

export { MessageContext, MessageProvider }
