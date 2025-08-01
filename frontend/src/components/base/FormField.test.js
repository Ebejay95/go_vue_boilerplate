import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import FormField from '@/components/base/FormField.vue'

describe('FormField Component', () => {
  const createWrapper = (props = {}) => {
    return mount(FormField, {
      props: {
        field: {
          key: 'test',
          label: 'Test Field',
          type: 'string',
          ...props.field
        },
        modelValue: '',
        ...props
      }
    })
  }

  describe('Text Inputs', () => {
    it('renders text input correctly', () => {
      const wrapper = createWrapper({
        field: { key: 'name', label: 'Name', type: 'string' }
      })

      expect(wrapper.find('input[type="text"]').exists()).toBe(true)
      expect(wrapper.find('label').text()).toBe('Name')
    })

    it('renders email input for email fields', () => {
      const wrapper = createWrapper({
        field: { key: 'email', label: 'Email' }
      })

      expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    })

    it('updates model value on input', async () => {
      const wrapper = createWrapper()
      const input = wrapper.find('input')

      await input.setValue('test value')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual(['test value'])
    })
  })

  describe('Select Fields', () => {
    it('renders select with options', () => {
      const wrapper = createWrapper({
        field: {
          key: 'role',
          inputType: 'select',
          options: [
            { value: 'user', label: 'User' },
            { value: 'admin', label: 'Admin' }
          ]
        }
      })

      const select = wrapper.find('select')
      expect(select.exists()).toBe(true)
      expect(wrapper.findAll('option')).toHaveLength(2)
    })
  })

  describe('Validation', () => {
    it('shows error message when error prop is provided', () => {
      const wrapper = createWrapper({
        error: 'This field is required',
        touched: true
      })

      expect(wrapper.find('.text-red-600').text()).toBe('This field is required')
    })

    it('shows required indicator', () => {
      const wrapper = createWrapper({
        field: { key: 'name', label: 'Name', required: true }
      })

      expect(wrapper.find('.text-red-500').text()).toBe('*')
    })
  })
})
