import React, { useCallback } from 'react'

import LinkIcon from '../../common/LinkIcon'
import { connectRouter, IRouterStateProps, IRouterDispatchProps } from '../../containers/RouterContainer'
import useKeyboard from '../../hooks/useKeyboard'
import DropdownMenu from '../../common/DropdownMenu'

type Props = {
  refresh: () => void
  markAllAsRead?: () => void
  canToggleStatus?: boolean
}

type AllProps = Props & IRouterStateProps & IRouterDispatchProps

function toggleSortOrderQueryParam(qs: string) {
  const params = new URLSearchParams(qs)
  params.set('sort', params.get('sort') === 'desc' ? 'asc' : 'desc')
  return params.toString()
}

function toggleStatusQueryParam(qs: string) {
  const params = new URLSearchParams(qs)
  params.set('status', params.get('status') === 'read' ? 'unread' : 'read')
  return params.toString()
}

export const ArticlesPageMenu = (props: AllProps) => {
  const {refresh, markAllAsRead, canToggleStatus, router, push} = props
  let { location: loc } = router

  const toggleSortOrder = useCallback((event: KeyboardEvent) => {
    event.preventDefault()
    push({...loc, search: toggleSortOrderQueryParam(loc.search)})
    return false
  }, [loc])
  
  const toggleStatus = useCallback((event: KeyboardEvent) => {
    if (!canToggleStatus) {
      return false
    }
    event.preventDefault()
    push({...loc, search: toggleStatusQueryParam(loc.search)})
    return false
  }, [loc])
  
  useKeyboard('shift+h', toggleStatus)
  useKeyboard('shift+o', toggleSortOrder)
  useKeyboard('shift+r', () => refresh())
  useKeyboard('shift+m', () => markAllAsRead ? markAllAsRead() : false)

  return (
    <DropdownMenu>
      <ul>
        <li>
          <LinkIcon onClick={refresh} icon="refresh">
            <span>Refresh</span><small>[shift+r]</small>
          </LinkIcon>
        </li>
        <li>
          <LinkIcon to={{...loc, search: toggleSortOrderQueryParam(loc.search)}} icon="sort">
            <span>Invert sort order</span><small>[shift+o]</small>
          </LinkIcon>
        </li>
        { markAllAsRead && <li>
          <LinkIcon onClick={markAllAsRead} icon="done_all">
            <span>Mark all as read</span><small>[shift+m]</small>
          </LinkIcon>
        </li> }
        { canToggleStatus && <li>
          <LinkIcon to={{...loc, search: toggleStatusQueryParam(loc.search)}} icon="history">
            <span>Toggle history</span><small>[shift+h]</small>
          </LinkIcon>
        </li> }
      </ul>
    </DropdownMenu>
  )
}

export default connectRouter(ArticlesPageMenu)
