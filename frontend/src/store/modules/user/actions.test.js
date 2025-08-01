import { describe, it, expect, vi, beforeEach } from 'vitest'
import actions from '@/store/modules/user/actions.js'

describe('User Store Actions', () => {
  let context
  let mockRootGetters
  let mockDispatch

  beforeEach(() => {
    mockDispatch = vi.fn()
    mockRootGetters = {
      'connection/isConnected': true,
      'connection/grpcClient': {
        listUsers: vi.fn(),
        createUser: vi.fn(),
        deleteUser: vi.fn()
      }
    }

    context = {
      commit: vi.fn(),
      dispatch: mockDispatch,
      rootGetters: mockRootGetters
    }
  })

  describe('fetchUsers', () => {
    it('fetches users successfully', async () => {
      const mockResponse = {
        getUsersList: () => [
          {
            getId: () => 1,
            getName: () => 'John Doe',
            getEmail: () => 'john@example.com',
            getAge: () => 30,
            getRole: () => 'user'
          }
        ]
      }

      mockDispatch.mockResolvedValue(mockResponse)

      const result = await actions.fetchUsers(context)

      expect(context.commit).toHaveBeenCalledWith('SET_LOADING', true)
      expect(context.commit).toHaveBeenCalledWith('SET_USERS', [
        {
          id: 1,
          name: 'John Doe',
          email: 'john@example.com',
          age: 30,
          role: 'user'
        }
      ])
      expect(result.count).toBe(1)
    })

    it('handles connection error', async () => {
      mockRootGetters['connection/isConnected'] = false

      await expect(actions.fetchUsers(context)).rejects.toThrow(
        'No gRPC connection available'
      )

      expect(context.commit).toHaveBeenCalledWith(
        'SET_ERROR',
        expect.stringContaining('No gRPC connection available')
      )
    })
  })

  describe('createUser', () => {
    it('creates user successfully', async () => {
      const userData = {
        name: 'Jane Doe',
        email: 'jane@example.com',
        age: 25,
        role: 'user'
      }

      const mockResponse = {
        getUser: () => ({
          getId: () => 2,
          getName: () => 'Jane Doe',
          getEmail: () => 'jane@example.com',
          getAge: () => 25,
          getRole: () => 'user'
        })
      }

      mockDispatch.mockResolvedValue(mockResponse)

      const result = await actions.createUser(context, userData)

      expect(context.commit).toHaveBeenCalledWith('SET_CREATING', true)
      expect(context.commit).toHaveBeenCalledWith('ADD_USER', {
        id: 2,
        name: 'Jane Doe',
        email: 'jane@example.com',
        age: 25,
        role: 'user'
      })
      expect(result.user.name).toBe('Jane Doe')
    })
  })
})
