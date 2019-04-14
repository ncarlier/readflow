import React, { useCallback } from 'react'

import LinkIcon from '../../common/LinkIcon'
import { connectRouter, IRouterStateProps, IRouterDispatchProps } from '../../containers/RouterContainer'
import useKeyboard from '../../hooks/useKeyboard'
import DropdownMenu from '../../common/DropdownMenu'
import ConfirmDialog from '../../common/ConfirmDialog';
import { useModal } from 'react-modal-hook';

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
    event.preventDefault()
    push({...loc, search: toggleStatusQueryParam(loc.search)})
    return false
  }, [loc])

  const [showMarkAllAsReadDialog, hideMarkAllAsReadDialog] = useModal(
    () => (
      <ConfirmDialog
        title="Mark all as read?"
        confirmLabel="Do it"
        onConfirm={() => {markAllAsRead!(); hideMarkAllAsReadDialog()}}
        onCancel={hideMarkAllAsReadDialog}
      >
        Are you sure to mark ALL articles as read?
      </ConfirmDialog>
    )
  )
  
  useKeyboard('shift+h', toggleStatus, canToggleStatus)
  useKeyboard('shift+o', toggleSortOrder)
  useKeyboard('shift+r', () => refresh())
  useKeyboard('shift+m', () => showMarkAllAsReadDialog(), markAllAsRead !== undefined)

  return (
    <DropdownMenu>
      <ul>
        <li>
          <LinkIcon onClick={refresh} icon="refresh">
            <span>Refresh</span><kbd>shift+r</kbd>
          </LinkIcon>
        </li>
        <li>
          <LinkIcon to={{...loc, search: toggleSortOrderQueryParam(loc.search)}} icon="sort">
            <span>Invert sort order</span><kbd>shift+o</kbd>
          </LinkIcon>
        </li>
        { markAllAsRead && <li>
          <LinkIcon onClick={showMarkAllAsReadDialog} icon="done_all">
            <span>Mark all as read</span><kbd>shift+m</kbd>
          </LinkIcon>
        </li> }
        { canToggleStatus && <li>
          <LinkIcon to={{...loc, search: toggleStatusQueryParam(loc.search)}} icon="history">
            <span>Toggle history</span><kbd>shift+h</kbd>
          </LinkIcon>
        </li> }
      </ul>
    </DropdownMenu>
  )
}

export default connectRouter(ArticlesPageMenu)
