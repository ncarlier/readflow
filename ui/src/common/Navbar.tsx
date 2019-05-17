import React from 'react'
import { useQuery } from 'react-apollo-hooks'

import { GetCategoriesResponse } from '../categories/models'
import { GetCategories } from '../categories/queries'
import useOnlineStatus from '../hooks/useOnlineStatus'
import { matchResponse } from './helpers'
import LinkIcon from './LinkIcon'
import Loader from './Loader'
import styles from './Navbar.module.css'
import Offline from './Offline'
import UserInfos from './UserInfos'

export default () => {
  const isOnline = useOnlineStatus()
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)

  const renderCategories = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader />,
    Error: err => <span>{err.message}</span>,
    Data: data => (
      <ul>
        {data.categories
          .filter(c => c.id !== null)
          .map(category => (
            <li key={`cat-${category.id}`}>
              <LinkIcon
                to={`/categories/${category.id}`}
                icon="bookmark"
                badge={category.unread ? category.unread : undefined}
              >
                {category.title}
              </LinkIcon>
            </li>
          ))}
      </ul>
    ),
    Other: () => <span>Unable to fetch categories!</span>
  })

  let total: number | undefined
  if (data && data.categories) {
    const all = data.categories.find(c => c.title === '_all')
    if (all) {
      total = all.unread
    }
  }

  return (
    <nav id="navbar" className={styles.nav}>
      <ul>
        <li>
          <h1>
            <img src={process.env.PUBLIC_URL + '/logo_white.svg'} />
          </h1>
          {isOnline && <UserInfos />}
          {!isOnline && <Offline />}
        </li>
        <li className={styles.links}>
          <span>Articles</span>
          <ul>
            <li>
              <LinkIcon to="/unread" icon="view_list" badge={total}>
                Articles to read
              </LinkIcon>
            </li>
            <li>
              <LinkIcon to="/offline" icon="signal_wifi_off">
                Offline articles
              </LinkIcon>
            </li>
            <li>
              <LinkIcon to="/history" icon="history">
                History
              </LinkIcon>
            </li>
          </ul>
        </li>
        <li className={styles.links}>
          <span>Categories</span>
          {isOnline && renderCategories(data, error, loading)}
        </li>
        <li className={styles.links}>
          <ul>
            <li>
              <LinkIcon id="navbar-link-settings" to="/settings" icon="settings">
                Settings
              </LinkIcon>
            </li>
          </ul>
        </li>
      </ul>
    </nav>
  )
}
