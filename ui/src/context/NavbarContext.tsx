import React, { createContext, ReactNode, useState } from 'react'

interface NavbarContextType {
  opened: boolean
  open: () => void
  close: () => void
}

const NavbarContext = createContext<NavbarContextType>({
  opened: true,
  open: () => null,
  close: () => null,
})

interface Props {
  children: ReactNode
}

const NavbarProvider = ({ children }: Props) => {
  const [opened, setOpened] = useState<boolean>(window.innerWidth > 767)

  const open = () => setOpened(true)
  const close = () => setOpened(false)

  return <NavbarContext.Provider value={{ opened, open, close }}>{children}</NavbarContext.Provider>
}

export { NavbarContext, NavbarProvider }
