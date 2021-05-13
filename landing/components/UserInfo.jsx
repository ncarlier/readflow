import useTranslation from 'next-translate/useTranslation'

const UserInfo = ({userData}) => {
  const { t } = useTranslation('common')
  if (!userData) {
    return <p>YOU ARE NOT AUTHENTICATED</p>
  }
  const {name, email} = userData.profile
  return (
    <section>
      <dl>
        <dt>Name</dt>
        <dd>{name}</dd>
        <dt>Email</dt>
        <dd>{email}</dd>
      </dl>
    </section>
  )
}

export default UserInfo
