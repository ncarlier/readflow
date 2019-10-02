export function getGQLError(err: any) {
  // console.log('Error', JSON.stringify(err, null, 4))
  if (err.networkError && err.networkError.name === 'ServerError') {
    return err.networkError.result.errors[0].message
  }
  return err.message
}
