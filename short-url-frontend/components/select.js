export default function Select({className, title, selectName, options, setOption, defaultValue}) {
  const handleSelectChange = (e) => setOption(e.target.value)

  return (
    <div className={className}>
      <label htmlFor={title} className="block text-md font-medium text-gray-700">
        {title}
        <select
          id={selectName}
          name={selectName}
          className="focus:ring-indigo-500 focus:border-indigo-500 h-full py-0 pl-2 pr-7 border-transparent bg-transparent text-gray-500 sm:text-sm rounded-md"
          value={defaultValue}
          onChange={e => handleSelectChange(e)}
        >
          { Object.entries(options).map(([key, value], i) => <option key={i} value={value}>{key}</option>) }
        </select>
      </label>
    </div>
  )
}