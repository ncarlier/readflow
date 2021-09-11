import React, { createContext, FC, useContext, useState } from 'react'

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

const NavbarProvider: FC = ({ children }) => {
  const [opened, setOpened] = useState<boolean>(window.innerWidth > 767)

  const open = () => setOpened(true)
  const close = () => setOpened(false)

  return <NavbarContext.Provider value={{ opened, open, close }}>{children}</NavbarContext.Provider>
}

export { NavbarProvider }

export const useNavbar = () => useContext(NavbarContext)
