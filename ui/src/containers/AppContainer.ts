import { connect } from 'react-redux'

import { AppState } from '../appStore'
import { ApplicationState } from '../store'

const mapStateToProps = ({ app }: ApplicationState): AppState => ({
  updateAvailable: app.updateAvailable,
  registration: app.registration,
})

export const connectApp = connect(mapStateToProps)
