import React from 'react'
import { useQuery } from '@apollo/client'
import { RouteComponentProps, withRouter } from 'react-router'
import { Link } from 'react-router-dom'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import { Loader, LinkIcon, NetworkStatus, Offline, UserInfos } from '../components'
import { matchResponse } from '../helpers'
import { ReactComponent as Logo } from './logo_header.min.svg'
import styles from './Navbar.module.css'
import { AddArticleLink } from '../articles/components'
import { Article } from '../articles/models'
import { useNavbar } from '../contexts'

export const Navbar = withRouter(({ location, history }: RouteComponentProps) => {
  const { pathname } = location
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)
  const navbar = useNavbar()

  const isCategoryActive = (id?: number) => {
    const _path = `/categories/${id}`
    return !!id && (pathname === _path || pathname.startsWith(`${_path}/`))
  }

  const menuAutoClose = () => {
    if (window.innerWidth <= 767) {
      setTimeout(navbar.close, 300)
    }
  }

  const redirectToNewArticle = (article: Article) => {
    history.push(`/inbox/${article.id}`)
    menuAutoClose()
  }

  const renderCategories = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader center />,
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
                badge={category.inbox}
              >
                {category.title}
              </LinkIcon>
            </li>
          ))}
      </ul>
    ),
  })

  return (
    <nav id="navbar" className={styles.nav}>
      <ul>
        <li>
          <h1>
            <Logo />
          </h1>
          <UserInfos />
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
                  to="/inbox"
                  icon="inbox"
                  badge={data && data.categories && data.categories._inbox}
                  active={pathname.startsWith('/inbox')}
                  onClick={menuAutoClose}
                >
                  Inbox
                </LinkIcon>
              </li>
              <li>
                <LinkIcon
                  as={Link}
                  to="/to_read"
                  icon="book"
                  badge={data && data.categories && data.categories._to_read}
                  active={pathname.startsWith('/to_read')}
                  onClick={menuAutoClose}
                >
                  To read
                </LinkIcon>
              </li>
              <li>
                <LinkIcon
                  as={Link}
                  to="/starred"
                  icon="star"
                  badge={data && data.categories && data.categories._starred}
                  active={pathname.startsWith('/starred')}
                  onClick={menuAutoClose}
                >
                  Starred
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
                Offline
              </LinkIcon>
            </li>
            <NetworkStatus status="online">
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
              <li>
                <AddArticleLink onSuccess={redirectToNewArticle} />
              </li>
            </NetworkStatus>
          </ul>
        </li>
        <NetworkStatus status="online">
          <li className={styles.links}>
            <span>Categories</span>
            {renderCategories(loading, data, error)}
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
