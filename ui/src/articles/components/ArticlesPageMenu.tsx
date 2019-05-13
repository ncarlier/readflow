import React, { useCallback } from 'react'
import { useModal } from 'react-modal-hook'

import ConfirmDialog from '../../common/ConfirmDialog'
import DropdownMenu from '../../common/DropdownMenu'
import Kbd from '../../common/Kbd'
import LinkIcon from '../../common/LinkIcon'
import { connectRouter, IRouterDispatchProps, IRouterStateProps } from '../../containers/RouterContainer'
import { DisplayMode } from './ArticlesDisplayMode'

interface Props {
  refresh: () => void
  markAllAsRead?: () => void
  mode: DisplayMode
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
  const { refresh, markAllAsRead, mode, router, push } = props
  let { location: loc } = router

  const toggleSortOrder = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      push({ ...loc, search: toggleSortOrderQueryParam(loc.search) })
      return false
    },
    [loc]
  )

  const toggleStatus = useCallback(
    (event: KeyboardEvent) => {
      event.preventDefault()
      push({ ...loc, search: toggleStatusQueryParam(loc.search) })
      return false
    },
    [loc]
  )

  const [showMarkAllAsReadDialog, hideMarkAllAsReadDialog] = useModal(() => (
    <ConfirmDialog
      title="Mark all as read?"
      confirmLabel="Do it"
      onConfirm={() => {
        if (markAllAsRead) markAllAsRead()
        hideMarkAllAsReadDialog()
      }}
      onCancel={hideMarkAllAsReadDialog}
    >
      Are you sure to mark ALL articles as read?
    </ConfirmDialog>
  ))

  return (
    <DropdownMenu>
      <ul>
        <li>
          <LinkIcon onClick={refresh} icon="refresh">
            <span>Refresh</span>
            <Kbd keys="shift+r" onKeypress={refresh} />
          </LinkIcon>
        </li>
        <li>
          <LinkIcon to={{ ...loc, search: toggleSortOrderQueryParam(loc.search) }} icon="sort">
            <span>Invert sort order</span>
            <Kbd keys="shift+o" onKeypress={toggleSortOrder} />
          </LinkIcon>
        </li>
        {markAllAsRead && (
          <li>
            <LinkIcon onClick={showMarkAllAsReadDialog} icon="done_all">
              <span>Mark all as read</span>
              <Kbd keys="shift+m" onKeypress={showMarkAllAsReadDialog} />
            </LinkIcon>
          </li>
        )}
        {mode === DisplayMode.category && (
          <li>
            <LinkIcon to={{ ...loc, search: toggleStatusQueryParam(loc.search) }} icon="history">
              <span>Toggle history</span>
              <Kbd keys="shift+h" onKeypress={toggleStatus} />
            </LinkIcon>
          </li>
        )}
      </ul>
    </DropdownMenu>
  )
}

export default connectRouter(ArticlesPageMenu)
