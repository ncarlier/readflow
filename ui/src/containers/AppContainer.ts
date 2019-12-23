import { connect } from 'react-redux'

import { ApplicationState } from '../store'
import { AppState } from '../appStore'

const mapStateToProps = ({ app }: ApplicationState): AppState => ({
  updateAvailable: app.updateAvailable
})

export const connectApp = connect(mapStateToProps)
