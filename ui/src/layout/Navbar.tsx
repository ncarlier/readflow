import React, { useContext } from 'react'
import { useQuery } from 'react-apollo-hooks'
import { RouteComponentProps, withRouter } from 'react-router'
import { Link } from 'react-router-dom'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import LinkIcon from '../components/LinkIcon'
import Loader from '../components/Loader'
import NetworkStatus from '../components/NetworkStatus'
import Offline from '../components/Offline'
import UserInfos from '../components/UserInfos'
import { NavbarContext } from '../context/NavbarContext'
import { matchResponse } from '../helpers'
import logo from './logo_header.svg'
import styles from './Navbar.module.css'

export default withRouter(({ location }: RouteComponentProps) => {
  const { pathname } = location
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)
  const navbar = useContext(NavbarContext)

  const isCategoryActive = (id?: number) => {
    const _path = `/categories/${id}`
    return !!id && (pathname === _path || pathname.startsWith(`${_path}/`))
  }

  const menuAutoClose = () => {
    if (window.innerWidth <= 767) {
      setTimeout(navbar.close, 300)
    }
  }

  const renderCategories = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader />,
    Error: (err) => <span>{err.message}</span>,
    Data: (data) => (
      <ul>
        {data.categories &&
          data.categories.entries.map((category) => (
            <li key={`cat-${category.id}`}>
              <LinkIcon
                as={Link}
                to={`/categories/${category.id}`}
                active={isCategoryActive(category.id)}
                icon="bookmark"
                onClick={menuAutoClose}
                badge={category.unread}
              >
                {category.title}
              </LinkIcon>
            </li>
          ))}
      </ul>
    ),
    Other: () => <span>Unable to fetch categories!</span>,
  })

  return (
    <nav id="navbar" className={styles.nav}>
      <ul>
        <li>
          <h1>
            <img src={logo} alt="readflow" />
          </h1>
          <NetworkStatus status="online">
            <UserInfos />
          </NetworkStatus>
          <NetworkStatus status="offline">
            <Offline />
          </NetworkStatus>
        </li>
        <li className={styles.links}>
          <span>Articles</span>
          <ul>
            <NetworkStatus status="online">
              <li>
                <LinkIcon
                  as={Link}
                  to="/unread"
                  icon="menu_book"
                  badge={data && data.categories && data.categories._all}
                  active={pathname.startsWith('/unread')}
                  onClick={menuAutoClose}
                >
                  Articles to read
                </LinkIcon>
              </li>
            </NetworkStatus>
            <li>
              <LinkIcon
                as={Link}
                to="/offline"
                icon="signal_wifi_off"
                active={pathname.startsWith('/offline')}
                onClick={menuAutoClose}
              >
                Offline articles
              </LinkIcon>
            </li>
            <NetworkStatus status="online">
              <li>
                <LinkIcon
                  as={Link}
                  to="/starred"
                  icon="star"
                  badge={data && data.categories && data.categories._starred}
                  active={pathname.startsWith('/starred')}
                  onClick={menuAutoClose}
                >
                  Starred articles
                </LinkIcon>
              </li>
              <li>
                <LinkIcon
                  as={Link}
                  to="/history"
                  icon="history"
                  active={pathname.startsWith('/history')}
                  onClick={menuAutoClose}
                >
                  History
                </LinkIcon>
              </li>
            </NetworkStatus>
          </ul>
        </li>
        <NetworkStatus status="online">
          <li className={styles.links}>
            <span>Categories</span>
            {renderCategories(data, error, loading)}
          </li>
          <li className={styles.links}>
            <ul>
              <li>
                <LinkIcon
                  id="navbar-link-settings"
                  as={Link}
                  to="/settings"
                  icon="settings"
                  onClick={menuAutoClose}
                  active={pathname.startsWith('/settings')}
                >
                  Settings
                </LinkIcon>
              </li>
            </ul>
          </li>
        </NetworkStatus>
      </ul>
    </nav>
  )
})
