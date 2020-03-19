import React from 'react'
import ReactDOM from 'react-dom'

import AboutButton from './AboutButton'

it('renders without crashing', () => {
  const div = document.createElement('div')
  ReactDOM.render(<AboutButton />, div)
  ReactDOM.unmountComponentAtNode(div)
})
