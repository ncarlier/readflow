export interface ArchiveService {
  id?: number
  alias: string
  provider: string
  config: string
  is_default: boolean
  created_at?: string
  updated_at?: string
}

export interface GetArchiveServicesResponse {
  archivers: ArchiveService[]
}

export interface GetArchiveServiceResponse {
  archiver: ArchiveService
}

export interface CreateOrUpdateArchiveServiceResponse {
  createOrUpdateArchiver: ArchiveService
}

export interface DeleteArchiveServiceRequest {
  ids: number[]
}

export interface DeleteArchiveServiceResponse {
  deleteArchivers: number
}
