

export type ArchiveService = {
  id?: number
  alias: string
  provider: string
  config: string
  is_default: boolean
  created_at?: string
  updated_at?: string
}

export type GetArchiveServicesResponse = {
  archivers: ArchiveService[]
}

export interface GetArchiveServiceResponse {
  archiver: ArchiveService
}

export type CreateOrUpdateArchiveServiceResponse = {
  createOrUpdateArchiver: ArchiveService
}
