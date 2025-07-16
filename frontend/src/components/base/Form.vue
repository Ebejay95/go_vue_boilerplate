// src/components/base/Form.vue
<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <div
      v-for="field in formFields"
      :key="field.key"
      class="form-field"
    >
      <!-- Dynamic Field Component -->
      <component
        :is="getFieldComponent(field)"
        v-bind="getFieldProps(field)"
        v-model="formData[field.key]"
        @update:modelValue="updateField(field.key, $event)"
        @blur="validateField(field.key)"
      />

      <!-- Field Error Display -->
      <div
        v-if="errors[field.key]"
        class="mt-1 text-sm text-red-600"
      >
        {{ errors[field.key] }}
      </div>
    </div>

    <!-- Form Actions -->
    <div class="flex space-x-3">
      <button
        type="submit"
        :disabled="isSubmitting || !isValid"
        class="btn primary"
      >
        <svg
          v-if="isSubmitting"
          class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        {{ isSubmitting ? submitLoadingText : submitText }}
      </button>

      <button
        v-if="showReset"
        type="button"
        @click="resetForm"
        class="btn flat"
      >
        {{ resetText }}
      </button>
    </div>
  </form>
</template>

<script>
import { ref, reactive, computed, watch } from 'vue'
import FormField from './FormField.vue'

export default {
  name: 'Form',

  components: {
    FormField
  },

  props: {
    // Proto message class (z.B. CreateUserRequest)
    protoMessage: {
      type: Function,
      required: true
    },

    // Field configurations override
    fieldConfig: {
      type: Object,
      default: () => ({})
    },

    // Initial data
    initialData: {
      type: Object,
      default: () => ({})
    },

    // Form behavior
    submitText: {
      type: String,
      default: 'Submit'
    },

    submitLoadingText: {
      type: String,
      default: 'Submitting...'
    },

    resetText: {
      type: String,
      default: 'Reset'
    },

    showReset: {
      type: Boolean,
      default: true
    },

    // Validation
    customValidators: {
      type: Object,
      default: () => ({})
    }
  },

  setup(props, { emit }) {
    // Reactive state using Vue 3 Composition API
    const formData = reactive({})
    const errors = reactive({})
    const touched = reactive({})
    const isSubmitting = ref(false)

    // Computed properties
    const formFields = computed(() => {
      return extractFieldsFromProto()
    })

    const isValid = computed(() => {
      return Object.keys(errors).length === 0
    })

    // Methods
    function extractFieldsFromProto() {
      // Create a temporary instance to inspect the proto structure
      const tempInstance = new props.protoMessage()
      const fields = []

      // Get all getters from the proto message
      const protoMethods = Object.getOwnPropertyNames(props.protoMessage.prototype)

      protoMethods.forEach(method => {
        if (method.startsWith('get') && method !== 'getClass') {
          const fieldName = getFieldNameFromGetter(method)
          const fieldType = inferFieldType(tempInstance, method)

          const fieldConfig = props.fieldConfig[fieldName] || {}

          fields.push({
            key: fieldName,
            getter: method,
            setter: method.replace('get', 'set'),
            type: fieldType,
            ...fieldConfig
          })
        }
      })

      return fields
    }

    function getFieldNameFromGetter(getter) {
      // Convert getName -> name, getEmail -> email, etc.
      return getter.replace('get', '').toLowerCase()
    }

    function inferFieldType(instance, getter) {
      try {
        const defaultValue = instance[getter]()

        if (typeof defaultValue === 'string') return 'string'
        if (typeof defaultValue === 'number') return 'number'
        if (typeof defaultValue === 'boolean') return 'boolean'
        if (Array.isArray(defaultValue)) return 'array'

        return 'string' // fallback
      } catch (error) {
        return 'string'
      }
    }

    function getFieldComponent(field) {
      // Return the appropriate component based on field configuration
      if (field.component) {
        return field.component
      }

      // Default components based on field type and configuration
      if (field.inputType === 'select') return 'FormField'
      if (field.inputType === 'textarea') return 'FormField'
      if (field.inputType === 'checkbox') return 'FormField'
      if (field.inputType === 'radio') return 'FormField'

      return 'FormField'
    }

    function getFieldProps(field) {
      const baseProps = {
        field: field,
        error: errors[field.key],
        touched: touched[field.key]
      }

      // Merge with custom props from field config
      return { ...baseProps, ...field.props }
    }

    function updateField(fieldKey, value) {
      formData[fieldKey] = value
      touched[fieldKey] = true

      // Clear error when user starts typing
      if (errors[fieldKey]) {
        delete errors[fieldKey]
      }
    }

    function validateField(fieldKey) {
      const field = formFields.value.find(f => f.key === fieldKey)
      const value = formData[fieldKey]

      // Required validation
      if (field.required && (!value || value === '')) {
        errors[fieldKey] = `${field.label || field.key} is required`
        return false
      }

      // Type validation
      if (value && !validateFieldType(field, value)) {
        errors[fieldKey] = `Invalid ${field.type} value`
        return false
      }

      // Custom validation
      if (field.validator && typeof field.validator === 'function') {
        const customError = field.validator(value, formData)
        if (customError) {
          errors[fieldKey] = customError
          return false
        }
      }

      // Global custom validators
      if (props.customValidators[fieldKey]) {
        const customError = props.customValidators[fieldKey](value, formData)
        if (customError) {
          errors[fieldKey] = customError
          return false
        }
      }

      // Clear error if validation passes
      delete errors[fieldKey]
      return true
    }

    function validateFieldType(field, value) {
      switch (field.type) {
        case 'number':
          return !isNaN(value) && isFinite(value)
        case 'string':
          return typeof value === 'string'
        case 'boolean':
          return typeof value === 'boolean'
        default:
          return true
      }
    }

    function validateForm() {
      let isValid = true

      formFields.value.forEach(field => {
        if (!validateField(field.key)) {
          isValid = false
        }
      })

      return isValid
    }

    async function handleSubmit() {
      if (!validateForm()) {
        return
      }

      isSubmitting.value = true

      try {
        // Create proto message instance
        const message = createProtoMessage()

        // Emit submit event with proto message
        emit('submit', message, formData)

      } catch (error) {
        console.dispatchError('Form submission error:', error)
        emit('error', error)
      } finally {
        isSubmitting.value = false
      }
    }

    function createProtoMessage() {
      const message = new props.protoMessage()

      formFields.value.forEach(field => {
        const value = formData[field.key]
        if (value !== undefined && value !== null && value !== '') {
          // Convert value to appropriate type
          let convertedValue = value

          if (field.type === 'number') {
            convertedValue = Number(value)
          } else if (field.type === 'boolean') {
            convertedValue = Boolean(value)
          }

          // Call the setter method
          message[field.setter](convertedValue)
        }
      })

      return message
    }

    function resetForm() {
      Object.keys(formData).forEach(key => {
        delete formData[key]
      })
      Object.keys(errors).forEach(key => {
        delete errors[key]
      })
      Object.keys(touched).forEach(key => {
        delete touched[key]
      })

      initializeForm()
      emit('reset')
    }

    function initializeForm() {
      // Initialize form data with initial values
      const initialFormData = { ...props.initialData }

      formFields.value.forEach(field => {
        if (initialFormData[field.key] === undefined) {
          initialFormData[field.key] = field.default || getDefaultValueForType(field.type)
        }
      })

      // Clear and set form data
      Object.keys(formData).forEach(key => {
        delete formData[key]
      })
      Object.assign(formData, initialFormData)
    }

    function getDefaultValueForType(type) {
      switch (type) {
        case 'string': return ''
        case 'number': return 0
        case 'boolean': return false
        case 'array': return []
        default: return ''
      }
    }

    // Watch for initial data changes
    watch(
      () => props.initialData,
      () => {
        initializeForm()
      },
      { deep: true }
    )

    // Initialize form on setup
    initializeForm()

    // Return reactive state and methods for template
    return {
      formData,
      errors,
      touched,
      isSubmitting,
      formFields,
      isValid,
      getFieldComponent,
      getFieldProps,
      updateField,
      validateField,
      handleSubmit,
      resetForm
    }
  }
}
</script>
