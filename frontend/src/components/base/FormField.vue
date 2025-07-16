// src/components/base/FormField.vue
<template>
  <div class="form-field-container">
    <!-- Field Label -->
    <label
      v-if="field.label !== false"
      :for="fieldId"
      class="block text-sm font-medium text-gray-700 mb-1"
    >
      {{ displayLabel }}
      <span v-if="field.required" class="text-red-500">*</span>
    </label>

    <!-- Text Input -->
    <input
      v-if="inputType === 'text' || inputType === 'email' || inputType === 'password' || inputType === 'url'"
      :id="fieldId"
      :type="inputType"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      @blur="$emit('blur')"
      :placeholder="field.placeholder || ''"
      :disabled="field.disabled"
      :readonly="field.readonly"
      :class="inputClasses"
    />

    <!-- Number Input -->
    <input
      v-else-if="inputType === 'number'"
      :id="fieldId"
      type="number"
      :value="modelValue"
      @input="$emit('update:modelValue', parseFloat($event.target.value) || 0)"
      @blur="$emit('blur')"
      :placeholder="field.placeholder || ''"
      :min="field.min"
      :max="field.max"
      :step="field.step || 1"
      :disabled="field.disabled"
      :readonly="field.readonly"
      :class="inputClasses"
    />

    <!-- Range Input -->
    <div v-else-if="inputType === 'range'" class="space-y-2">
      <input
        :id="fieldId"
        type="range"
        :value="modelValue"
        @input="$emit('update:modelValue', parseFloat($event.target.value))"
        @blur="$emit('blur')"
        :min="field.min || 0"
        :max="field.max || 100"
        :step="field.step || 1"
        :disabled="field.disabled"
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>{{ field.min || 0 }}</span>
        <span class="font-medium">{{ modelValue }}</span>
        <span>{{ field.max || 100 }}</span>
      </div>
    </div>

    <!-- Textarea -->
    <textarea
      v-else-if="inputType === 'textarea'"
      :id="fieldId"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      @blur="$emit('blur')"
      :placeholder="field.placeholder || ''"
      :rows="field.rows || 3"
      :disabled="field.disabled"
      :readonly="field.readonly"
      :class="inputClasses"
    ></textarea>

    <!-- Select -->
    <select
      v-else-if="inputType === 'select'"
      :id="fieldId"
      :value="modelValue"
      @change="$emit('update:modelValue', $event.target.value)"
      @blur="$emit('blur')"
      :disabled="field.disabled"
      :class="inputClasses"
    >
      <option v-if="field.placeholder" value="" disabled>
        {{ field.placeholder }}
      </option>
      <option
        v-for="option in field.options"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </option>
    </select>

    <!-- Checkbox -->
    <div v-else-if="inputType === 'checkbox'" class="flex items-center">
      <input
        :id="fieldId"
        type="checkbox"
        :checked="modelValue"
        @change="$emit('update:modelValue', $event.target.checked)"
        @blur="$emit('blur')"
        :disabled="field.disabled"
        class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
      />
      <label
        :for="fieldId"
        class="ml-2 block text-sm text-gray-900"
      >
        {{ field.checkboxLabel || displayLabel }}
      </label>
    </div>

    <!-- Radio Group -->
    <div v-else-if="inputType === 'radio'" class="space-y-2">
      <div
        v-for="option in field.options"
        :key="option.value"
        class="flex items-center"
      >
        <input
          :id="`${fieldId}-${option.value}`"
          type="radio"
          :value="option.value"
          :checked="modelValue === option.value"
          @change="$emit('update:modelValue', option.value)"
          @blur="$emit('blur')"
          :disabled="field.disabled"
          class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300"
        />
        <label
          :for="`${fieldId}-${option.value}`"
          class="ml-2 block text-sm text-gray-900"
        >
          {{ option.label }}
        </label>
      </div>
    </div>

    <!-- Date Input -->
    <input
      v-else-if="inputType === 'date'"
      :id="fieldId"
      type="date"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      @blur="$emit('blur')"
      :disabled="field.disabled"
      :readonly="field.readonly"
      :class="inputClasses"
    />

    <!-- Color Input -->
    <input
      v-else-if="inputType === 'color'"
      :id="fieldId"
      type="color"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      @blur="$emit('blur')"
      :disabled="field.disabled"
      :readonly="field.readonly"
      :class="inputClasses"
    />

    <!-- File Input -->
    <input
      v-else-if="inputType === 'file'"
      :id="fieldId"
      type="file"
      @change="handleFileChange"
      @blur="$emit('blur')"
      :multiple="field.multiple"
      :accept="field.accept"
      :disabled="field.disabled"
      :class="inputClasses"
    />

    <!-- Help Text -->
    <p
      v-if="field.help"
      class="mt-1 text-sm text-gray-500"
    >
      {{ field.help }}
    </p>

    <!-- Error Message -->
    <p
      v-if="error && touched"
      class="mt-1 text-sm text-red-600"
    >
      {{ error }}
    </p>
  </div>
</template>

<script>
export default {
  name: 'FormField',

  props: {
    field: {
      type: Object,
      required: true
    },
    modelValue: {
      type: [String, Number, Boolean, Array, Object],
      default: ''
    },
    error: {
      type: String,
      default: ''
    },
    touched: {
      type: Boolean,
      default: false
    }
  },

  emits: ['update:modelValue', 'blur'],

  computed: {
    fieldId() {
      return `field-${this.field.key}-${this._uid}`
    },

    displayLabel() {
      return this.field.label || this.capitalizeFirst(this.field.key)
    },

    inputType() {
      // Custom input type from field config
      if (this.field.inputType) {
        return this.field.inputType
      }

      // Auto-detect based on field name and type
      if (this.field.key.includes('email')) return 'email'
      if (this.field.key.includes('password')) return 'password'
      if (this.field.key.includes('url') || this.field.key.includes('website')) return 'url'
      if (this.field.key.includes('age') || this.field.key.includes('count')) return 'number'
      if (this.field.key.includes('date')) return 'date'
      if (this.field.key.includes('color')) return 'color'
      if (this.field.key.includes('file')) return 'file'

      // Based on proto type
      if (this.field.type === 'number') return 'number'
      if (this.field.type === 'boolean') return 'checkbox'

      return 'text'
    },

    inputClasses() {
      const baseClasses = [
        'block w-full rounded-md border-gray-300 shadow-sm',
        'focus:border-blue-500 focus:ring-blue-500',
        'disabled:bg-gray-50 disabled:cursor-not-allowed'
      ]

      if (this.error && this.touched) {
        baseClasses.push('border-red-300 focus:border-red-500 focus:ring-red-500')
      }

      return baseClasses.join(' ')
    }
  },

  methods: {
    capitalizeFirst(str) {
      return str.charAt(0).toUpperCase() + str.slice(1)
    },

    handleFileChange(event) {
      const files = event.target.files
      if (this.field.multiple) {
        this.$emit('update:modelValue', Array.from(files))
      } else {
        this.$emit('update:modelValue', files[0] || null)
      }
    }
  }
}
</script>

<style scoped>
.form-field-container {
  @apply w-full;
}

/* Range slider styling */
input[type="range"]::-webkit-slider-thumb {
  @apply appearance-none w-5 h-5 bg-blue-600 rounded-full cursor-pointer;
}

input[type="range"]::-moz-range-thumb {
  @apply w-5 h-5 bg-blue-600 rounded-full cursor-pointer border-0;
}
</style>