import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import { getUserInfo } from '@/api/users'
import store from '@/store'

export interface IUserState {
  name: string
  displayName: string
  profileImageURL: string
  admin: boolean
}

@Module({ dynamic: true, store, name: 'user' })
class User extends VuexModule implements IUserState {
  public name = ''
  public displayName = ''
  public profileImageURL = ''
  public admin: boolean = false

  @Mutation
  private SET_NAME(name: string) {
    this.name = name
  }

  @Mutation
  private SET_DISPLAY_NAME(displayName: string) {
    this.displayName = displayName
  }

  @Mutation
  private SET_PROFILE_IMAGE_URL(profileImageURL: string) {
    this.profileImageURL = profileImageURL
  }

  @Mutation
  private SET_ADMIN(admin: boolean) {
    this.admin = admin
  }

  @Action
  public async GetUserInfo() {
    const user = await getUserInfo()
    if (!user) {
      throw Error('Verification failed, please Login again.')
    }
    const { admin, name, displayName, profileImageURL } = user.data
    this.SET_ADMIN(admin)
    this.SET_NAME(name)
    this.SET_DISPLAY_NAME(displayName)
    this.SET_PROFILE_IMAGE_URL(profileImageURL)
  }

  @Action
  public async Login() {
    window.location.assign(process.env.VUE_APP_BASE_API + 'auth/github')
  }

  @Action
  public async Logout() {
    window.location.assign(process.env.VUE_APP_BASE_API + 'auth/logout')
  }
}

export const UserModule = getModule(User)
