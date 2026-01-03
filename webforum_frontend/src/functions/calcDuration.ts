export default function calcDuration(dateString: string){
        const now = new Date()
        const past = new Date(dateString)
        const diffMs = now.getTime() - past.getTime()   //in milliseconds
        const diffSec = Math.floor(diffMs/1000)
        const diffMin = Math.floor(diffSec/60)
        const diffHr = Math.floor(diffMin/60)
        const diffDay = Math.floor(diffHr/24)

        if (diffDay > 0) return `${diffDay} day${diffDay>1? "s" : ""} ago`
        if (diffHr > 0) return `${diffHr} hr${diffHr>1? "s" : ""} ago`
        if (diffMin > 0) return `${diffMin} min${diffMin>1? "s" : ""} ago`
        return "now"
    }