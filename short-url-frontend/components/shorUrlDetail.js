import Date from './date'

export default function ShortUrlDetail({ hash, data }) {
  return (
    <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
      <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
        <div className="flex flex-wrapbg-white py-8 px-10 rounded-md shadow-md max-w-lg mx-auto">
          <div className="flex flex-col">
            <div className="flex flex-row items-center">
              <div>
                <p className="text-3xl font-bold text-black mr-4">{hash}</p>
              </div>
              <div>
                {"owner" in data?
                  <span
                    className="px-2 inline-block items-end text-md font-semibold rounded-full bg-purple-100 text-purple-700"
                  >{data["owner"]}</span>
                  :
                  <></>
                }
              </div>
            </div>
            <div className="flex flex-row mb-2 items-center">
              <span
                className="px-2 inline-flex text-md leading-5 font-semibold rounded-full bg-green-100 text-green-800"
              >{data["type"]}</span>
              &nbsp;
              <h3 className="text-lg leading-6 font-medium text-green-600">{data["target"]}</h3>
            </div>
            <div className="flex flex-row mb-2 items-center">
              <span
                className="px-2 inline-flex text-md leading-5 font-semibold rounded-full bg-blue-100 text-blue-800"
              >Created at:&nbsp;<Date dateString={data["createdAt"]} /></span>
            </div>
            <div className="flex flex-row mb-2 items-center">
              <span
                className="px-2 inline-flex text-md leading-5 font-semibold rounded-full bg-yellow-100 text-yellow-800"
              >Redirect count:&nbsp;{data["count"]}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}